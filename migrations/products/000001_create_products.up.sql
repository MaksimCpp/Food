CREATE TABLE products(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    description VARCHAR(3000) NOT NULL,
    price INTEGER NOT NULL, -- цена в копейках
    image_url VARCHAR(512),
    is_available NOT NULL DEFAULT true
);