
CREATE TYPE user_role AS ENUM (
  'user',
  'admin'
);


CREATE TABLE users (
  id uuid PRIMARY KEY,
  user_role user_role NOT NULL,
  first_name varchar(50) NOT NULL,
  last_name varchar(50) NOT NULL,
  password varchar(150) NOT NULL,
  email varchar(50) UNIQUE NOT NULL,
  phone_number varchar(50),
  created_at timestamp NOT NULL DEFAULT 'now()',
  updated_at timestamp NOT NULL DEFAULT 'now()'
);


CREATE TABLE session (
  id uuid PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  user_agent text NOT NULL,
  ip_address varchar(64) NOT NULL,
  is_active bool NOT NULL,
  expires_at timestamp,
  last_active_at timestamp,
  created_at timestamp NOT NULL DEFAULT 'now()',
  updated_at timestamp NOT NULL DEFAULT 'now()'
);

CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at INT DEFAULT 0
);



CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    price NUMERIC(10, 2) NOT NULL CHECK (price >= 0),
    product_size varchar(50),
    color  varchar(50),
    sale_price  NUMERIC(10, 2),
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);


CREATE TABLE IF NOT EXISTS products_pictures (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    picture_url TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);  



-- Create basket_items table
CREATE TABLE basket_items (
    id UUID PRIMARY KEY,
    order_id varchar(50),
    product_id UUID NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    price NUMERIC(10, 2) NOT NULL,
    count INT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'not_sold' CHECK (status IN ('sold', 'not_sold')),
    deleted_at TIMESTAMP DEFAULT NULL
);






CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    total_price NUMERIC(10, 2) NOT NULL CHECK (total_price >= 0),
    status VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);





