package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"aura-fashion/config"
	"aura-fashion/internal/entity"
	"aura-fashion/pkg/etc"
	"aura-fashion/pkg/hash"
	"aura-fashion/pkg/jwt"

	"github.com/gin-gonic/gin"
)

// Login godoc
// @Router /auth/login [post]
// @Summary Login
// @Description Login
// @Tags auth
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Param body body entity.LoginRequest true "User"
// @Success 200 {object} entity.SuccessResponse
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) Login(ctx *gin.Context) {
	var (
		body entity.LoginRequest
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	user, err := h.UseCase.UserRepo.GetSingle(ctx, entity.UserSingleRequest{
		Email:    body.Email,
	})
	if h.HandleDbError(ctx, err, "Error getting user") {
		return
	}



	if !hash.CheckPasswordHash(body.Password, user.Password) {
		h.ReturnError(ctx, config.ErrorInvalidPass, "Incorrect password", http.StatusBadRequest)
		return
	}

	// create session
	newSession := entity.Session{
		UserID:       user.ID,
		IPAddress:    ctx.ClientIP(),
		ExpiresAt:    time.Now().Add(time.Hour * 999999).Format(time.RFC3339),
		UserAgent:    ctx.Request.UserAgent(),
		IsActive:     true,
		LastActiveAt: time.Now().Format(time.RFC3339),
	}

	session, err := h.UseCase.SessionRepo.Create(ctx, newSession)
	if h.HandleDbError(ctx, err, "Error while creating new session") {
		return
	}

	// generate jwt token
	jwtFields := map[string]interface{}{
		"sub":        user.ID,
		"user_role":  user.UserRole,
		"session_id": session.ID,
	}

	user.AccessToken, err = jwt.GenerateJWT(jwtFields, h.Config.JWT.Secret)
	if err != nil {
		h.ReturnError(ctx, config.ErrorInternalServer, "Oops, something went wrong!!!", http.StatusInternalServerError)
		return
	}

	ctx.JSON(200, user)
}

// // Logout godoc
// // @Router /auth/logout [post]
// // @Summary Logout
// // @Description Logout
// // @Security BearerAuth
// // @Tags auth
// // @Accept  json
// // @Produce  json
// // @Success 200 {object} entity.SuccessResponse
// // @Failure 400 {object} entity.ErrorResponse
// func (h *Handler) Logout(ctx *gin.Context) {
// 	sessionID := ctx.GetHeader("session_id")
// 	if sessionID == "" {
// 		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid session ID", 400)
// 		return
// 	}

// 	err := h.UseCase.SessionRepo.Delete(ctx, entity.Id{
// 		ID: sessionID,
// 	})
// 	if h.HandleDbError(ctx, err, "Error deleting session") {
// 		return
// 	}

// 	ctx.JSON(200, entity.SuccessResponse{
// 		Message: "Successfully logged out",
// 	})
// }

// Register godoc
// @Router /auth/register [post]
// @Summary Register
// @Description Register
// @Security BearerAuth
// @Tags auth
// @Accept  json
// @Produce  json
// @Param body body entity.RegisterRequest true "User"
// @Success 200 {object} entity.User
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) Register(ctx *gin.Context) {
	var body entity.User

	// Parse request body
	if err := ctx.ShouldBindJSON(&body); err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}

	// Check if user already exists
	_, err := h.UseCase.UserRepo.GetSingle(ctx, entity.UserSingleRequest{
		Email: body.Email,
	})
	if err == nil {
		h.ReturnError(ctx, config.ErrorConflict, "User already exists", 400)
		return
	}

	// Hash password
	body.Password, err = hash.HashPassword(body.Password)
	if err != nil {
		h.ReturnError(ctx, config.ErrorInternalServer, "Error hashing password", http.StatusInternalServerError)
		return
	}

	body.UserRole="user"
	userJSON, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	err = h.Redis.Set(ctx, "user:"+body.Email, string(userJSON),15*60)
	if err != nil {
		fmt.Println("Error storing user in Redis:", err)
		return
	}
	// Generate OTP
	otp := etc.GenerateOTP(6)

	// Set OTP in Redis and send email in parallel
	errChan := make(chan error, 2)

	// Redis operation
	go func() {
		errChan <- h.Redis.Set(ctx, fmt.Sprintf("otp-%s", body.Email), otp, 5*60)
	}()

	// Email operation
	go func() {
		emailBody, err := etc.GenerateOtpEmailBody(otp)
		if err != nil {
			errChan <- fmt.Errorf("error generating OTP email body: %w", err)
			return
		}
		errChan <- etc.SendEmail(h.Config.Gmail.Host, h.Config.Gmail.Port, h.Config.Gmail.Email, h.Config.Gmail.EmailPass, body.Email, emailBody)
	}()

	// Collect errors
	for i := 0; i < 2; i++ {
		if err := <-errChan; err != nil {
			h.ReturnError(ctx, config.ErrorInternalServer, err.Error(), 500)
			return
		}
	}

	// Success response
	ctx.JSON(201, entity.SuccessResponse{
		Message: "User registered successfully, please verify your email address",
	})
}

// VerifyEmail godoc
// @Router /auth/verify-email [post]
// @Summary VerifyEmail
// @Description VerifyEmail
// @Security BearerAuth
// @Tags auth
// @Accept  json
// @Produce  json
// @Param body body entity.VerifyEmail true "User"
// @Success 200 {object} entity.User
// @Failure 400 {object} entity.ErrorResponse
func (h *Handler) VerifyEmail(ctx *gin.Context) {
	var (
		body entity.VerifyEmail
	)

	err := ctx.ShouldBindJSON(&body)
	if err != nil {
		h.ReturnError(ctx, config.ErrorBadRequest, "Invalid request body", 400)
		return
	}
	
	key := fmt.Sprintf("otp-%s", body.Email)

	otp, err := h.Redis.Get(ctx, key)
	if err != nil {
		h.HandleDbError(ctx, err, "error from get redis ")
		return
	}

	if otp != body.Otp {
		h.ReturnError(ctx, config.ErrorBadRequest, "Incorrect otp", http.StatusBadRequest)
		return
	}

	// Redisdan foydalanuvchini olish
	val, err := h.Redis.Get(ctx, "user:"+body.Email)
	if err != nil {
		fmt.Println("User not found in Redis or expired")
		return
	}

	// JSON'dan struct'ga parse qilish
	var storedUser entity.User
	err = json.Unmarshal([]byte(val), &storedUser)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}


	user, err:=h.UseCase.UserRepo.Create(ctx,storedUser)
	if err!=nil{
		h.HandleDbError(ctx,err, "error while create user")
		return
	}

	// create session
	newSession := entity.Session{
		UserID:       user.ID,
		IPAddress:    ctx.ClientIP(),
		ExpiresAt:    time.Now().Add(time.Hour * 999999).Format(time.RFC3339),
		UserAgent:    ctx.Request.UserAgent(),
		IsActive:     true,
		LastActiveAt: time.Now().Format(time.RFC3339),
	}

	session, err := h.UseCase.SessionRepo.Create(ctx, newSession)
	if h.HandleDbError(ctx, err, "Error while creating new session") {
		return
	}
	
	// generate jwt token
	jwtFields := map[string]interface{}{
		"sub":        user.ID,
		"user_role":  user.UserRole,
		"session_id": session.ID,
	}

	user.AccessToken, err = jwt.GenerateJWT(jwtFields, h.Config.JWT.Secret)
	if err != nil {
		
		h.ReturnError(ctx, config.ErrorInternalServer, "Oops, something went wrong!!!", http.StatusInternalServerError)
		return
	}

	ctx.JSON(200, user)
}
