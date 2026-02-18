-- +goose Up
-- +goose StatementBegin
CREATE TABLE metrics (
    id VARCHAR(255) UNIQUE,
    type TEXT NOT NULL CHECK (type IN ('gauge','counter')),
    value DOUBLE PRECISION,
    delta BIGINT,
    CONSTRAINT match_type CHECK (
        (type='gauge' AND value IS NOT NULL AND delta IS NULL) OR
        (type='counter' AND delta IS NOT NULL AND value IS NULL)
        )
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE metrics;
-- +goose StatementEnd
