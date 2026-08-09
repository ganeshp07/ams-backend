CREATE TABLE IF NOT EXISTS courses (
    id              BIGSERIAL PRIMARY KEY,
    code            VARCHAR(20)  NOT NULL UNIQUE,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    credits         INT          NOT NULL DEFAULT 3,
    course_type     VARCHAR(50)  NOT NULL DEFAULT 'Theory',
    program_id      BIGINT       NOT NULL REFERENCES programs(id),
    semester_number INT          NOT NULL,
    is_active       BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
