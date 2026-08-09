CREATE TABLE IF NOT EXISTS programs (
    id             BIGSERIAL PRIMARY KEY,
    name           VARCHAR(255) NOT NULL,
    code           VARCHAR(50)  NOT NULL UNIQUE,
    description    TEXT,
    duration_years INT          NOT NULL DEFAULT 4,
    department_id  BIGINT       NOT NULL REFERENCES departments(id),
    is_active      BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
