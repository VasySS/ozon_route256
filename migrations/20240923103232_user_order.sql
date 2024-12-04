-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS user_order (
    id BIGINT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id BIGINT NOT NULL,
    weight REAL NOT NULL,
    price REAL NOT NULL,
    expiry_date TIMESTAMP NOT NULL,
    receive_date TIMESTAMP NOT NULL DEFAULT '0001-01-01 00:00:00'
);

CREATE INDEX user_order_user_id_idx ON user_order USING btree (user_id);

COMMENT ON TABLE user_order IS 'Заказы пользователя, принятые от курьера в ПВЗ';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_order;
-- +goose StatementEnd
