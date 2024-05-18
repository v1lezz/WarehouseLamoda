-- +goose Up
-- +goose StatementBegin
CREATE TABLE goods (
    name VARCHAR(255) NOT NULL,
    size INTEGER NOT NULL CHECK (size > 0),
    id INTEGER NOT NULL CHECK (id > 0) PRIMARY KEY
);

CREATE TABLE warehouse(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    is_available BOOLEAN NOT NULL
);

CREATE TABLE goods_warehouse(
    id SERIAL PRIMARY KEY,
    warehouse_id SERIAL NOT NULL REFERENCES warehouse(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    good_id INTEGER NOT NULL REFERENCES goods(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    count INTEGER NOT NULL CHECK (count > 0),
    reserved INTEGER NOT NULL,
    UNIQUE(warehouse_id, good_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE goods_warehouse;
DROP TABLE warehouse;
DROP TABLE goods;
-- +goose StatementEnd
