-- +goose Up
-- +goose StatementBegin
CREATE TABLE metrics (
    id VARCHAR(255) UNIQUE NOT NULL,
    type TEXT CHECK ( type IN('gauge', 'counter') ),
    value DOUBLE PRECISION,
    delta BIGINT
);

CREATE INDEX idx_metrics_id ON metrics(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE metrics;
-- +goose StatementEnd
