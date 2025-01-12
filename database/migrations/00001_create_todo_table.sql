-- +goose Up
-- +goose StatementBegin

CREATE TYPE todo_status AS ENUM ('pending', 'completed');

CREATE TABLE IF NOT EXISTS todo (
  id UUID PRIMARY KEY,
  title VARCHAR NOT NULL,
  description TEXT,
  status todo_status NOT NULL DEFAULT 'pending',
  due_date DATE NOT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS todo;
DROP TYPE IF EXISTS todo_status;
-- +goose StatementEnd
