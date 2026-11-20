CREATE TYPE order_status AS ENUM (
    'accepted',
    'out_for_delivery',
    'delivered'
);

CREATE TABLE orders(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL CHECK (user_id > 0),
    product_id BIGINT NOT NULL CHECK (product_id > 0),
    products_quantity INTEGER NOT NULL CHECK (products_quantity > 0),
    courier_id BIGINT CHECK (courier_id > 0),
    address VARCHAR(400) NOT NULL, -- в будущем добавить GEOGRAPHY
    status order_status NOT NULL DEFAULT 'accepted',

    CONSTRAINT orders_status_courier_check CHECK (
        (courier_id IS NULL AND status = 'accepted')
        OR
        (courier_id IS NOT NULL AND status IN ('out_for_delivery', 'delivered'))
    )
);