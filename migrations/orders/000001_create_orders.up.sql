CREATE TYPE order_status AS ENUM ();

CREATE TABLE orders(
    user_id INTEGER NOT NULL CHECK (user_id > 0),
    product_id INTEGER NOT NULL CHECK (product_id > 0),
    products_quantity INTEGER NOT NULL CHECK (products_quantity > 0),
    courier_id INTEGER CHECK (courier_id > 0),
    address VARCHAR(400) NOT NULL, -- в будущем добавить GEOGRAPHY
    status order_status NOT NULL DEFAULT '',
);