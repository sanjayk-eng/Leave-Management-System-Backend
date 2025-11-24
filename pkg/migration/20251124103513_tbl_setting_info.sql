-- +goose Up



CREATE TABLE IF NOT EXISTS Tbl_Company_Settings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    working_days_per_month INT NOT NULL DEFAULT 22,
    allow_manager_add_leave BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
