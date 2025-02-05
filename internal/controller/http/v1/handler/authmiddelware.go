package handler

import (
	"fmt"
	"net/http"
	"strings"

	"aura-fashion/internal/entity"
	"aura-fashion/pkg/jwt"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func (h *Handler) AuthMiddleware(e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			userRole string
			act      = c.Request.Method
			obj      = c.FullPath()
		)

		token := c.GetHeader("Authorization")
		if token == "" {
			userRole = "unauthorized"
		}

		if userRole == "" {
			token = strings.TrimPrefix(token, "Bearer ")

			claims, err := jwt.ParseJWT(token, h.Config.JWT.Secret)
			if err != nil {
				userRole = "unauthorized"
			}

			v, ok := claims["user_role"].(string)
			if !ok {
				userRole = "unauthorized"
			} else {
				userRole = v
			}

			for key, value := range claims {
				c.Request.Header.Set(key, fmt.Sprintf("%v", value))
			}
		}

		// TO DO: Check if session is valid

		if userRole != "unauthorized" {
			session, err := h.UseCase.SessionRepo.GetSingle(c, entity.Id{ID: c.GetHeader("session_id")})
			if err != nil {
				fmt.Println("error while gettign single session", err)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Session is invalid"})
				return
			}

			if !session.IsActive {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Session is not active"})
				return
			}
		}

		ok, err := e.EnforceSafe(userRole, obj, act)
		fmt.Println("role: ", userRole)
		fmt.Println("path: ", obj)
		fmt.Println("method: ", act)
		if err != nil {
			h.Logger.Error(err, "Error enforcing policy")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}

		c.Next()
	}
}

func (h *Handler) GetIdFromToken(c *gin.Context) (string, int) {
	var softToken string
	token := c.GetHeader("Authorization")
	if token == "" {
		return "unauthorized", http.StatusUnauthorized
	} else if strings.Contains(token, "Bearer") {
		softToken = strings.TrimPrefix(token, "Bearer ")
	} else {
		softToken = token
	}

	claims, err := jwt.ParseJWT(softToken, h.Config.JWT.Secret)
	if err != nil {
		return "unauthorized", http.StatusUnauthorized
	}

	return cast.ToString(claims["sub"]), 0
}
