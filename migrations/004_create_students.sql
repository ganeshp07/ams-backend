CREATE TABLE IF NOT EXISTS students (
    id               BIGSERIAL PRIMARY KEY,
    user_id          BIGINT       NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    roll_number      VARCHAR(50)  NOT NULL UNIQUE,
    first_name       VARCHAR(100) NOT NULL,
    last_name        VARCHAR(100) NOT NULL,
    phone            VARCHAR(20),
    date_of_birth    DATE,
    gender           VARCHAR(20),
    address          TEXT,
    program_id       BIGINT       NOT NULL REFERENCES programs(id),
    admission_year   INT          NOT NULL,
    current_semester INT          NOT NULL DEFAULT 1,
    is_active        BOOLEAN      NOT NULL DEFAULT TRUE,
    created_by       BIGINT       REFERENCES users(id),
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
