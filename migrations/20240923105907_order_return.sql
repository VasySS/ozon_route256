-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS order_return (
    id BIGINT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL,
    CONSTRAINT order_id_fk 
        FOREIGN KEY (order_id) 
        REFERENCES user_order(id)
        ON DELETE CASCADE
);

CREATE INDEX order_return_order_id_idx ON order_return USING btree (order_id);

COMMENT ON TABLE order_return IS 'Возвраты, принятые от пользователя в ПВЗ';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS order_return;
-- +goose StatementEnd
