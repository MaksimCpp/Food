CREATE TABLE products(
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(256) NOT NULL,
    description VARCHAR(3000) NOT NULL,
    price INTEGER NOT NULL CHECK (price >= 0), -- цена в копейках
    image_key VARCHAR(512),
    is_available BOOLEAN NOT NULL DEFAULT true,
    stock_quantity INTEGER NOT NULL DEFAULT 0
        CHECK (reserved_quantity >= 0),
    reserved_quantity INTEGER NOT NULL DEFAULT 0 
        CHECK (reserved_quantity <= stock_quantity),
        CHECK (reserved_quantity >= 0)
);