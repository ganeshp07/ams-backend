CREATE TABLE IF NOT EXISTS faculty (
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT       NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    employee_id   VARCHAR(50)  NOT NULL UNIQUE,
    first_name    VARCHAR(100) NOT NULL,
    last_name     VARCHAR(100) NOT NULL,
    phone         VARCHAR(20),
    department_id BIGINT       REFERENCES departments(id),
    designation   VARCHAR(100),
    is_active     BOOLEAN      NOT NULL DEFAULT TRUE,
    created_by    BIGINT       REFERENCES users(id),
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
