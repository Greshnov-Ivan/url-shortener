-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "links" (
    "id" BIGSERIAL PRIMARY KEY,
    "source_url" VARCHAR(2048) NOT NULL UNIQUE,
    "created_at" TIMESTAMPTZ NOT NULL,
    "expires_at" TIMESTAMPTZ NULL,
    "last_requested_at" TIMESTAMPTZ NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "links";
-- +goose StatementEnd
