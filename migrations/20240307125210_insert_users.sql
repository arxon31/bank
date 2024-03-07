-- +goose Up
-- +goose StatementBegin
INSERT INTO users (id, name, amount) VALUES
(1, 'Elizar', 9999999), (2, 'Vasya', 9999999), (3, 'Petya', 100), (4, 'Kolya', 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users;
-- +goose StatementEnd
