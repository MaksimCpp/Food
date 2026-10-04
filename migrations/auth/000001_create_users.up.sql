CREATE TYPE user_role AS ENUM ('user', 'admin', 'courier');

CREATE TABLE users(
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(256) NOT NULL,
    email VARCHAR(256) UNIQUE NOT NULL,
    password_hash VARCHAR(300) NOT NULL,
    role user_role NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT NOW()
);