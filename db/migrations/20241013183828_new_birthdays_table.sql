-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS birthdays(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name varchar(255) NOT NULL,
    date DATE NOT NULL,
    additional TEXT,
    user_id BIGINT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS birthdays;
-- +goose StatementEnd
