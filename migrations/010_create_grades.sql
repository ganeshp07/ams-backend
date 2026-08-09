CREATE TABLE IF NOT EXISTS grades (
    id              BIGSERIAL PRIMARY KEY,
    registration_id BIGINT      NOT NULL UNIQUE REFERENCES registrations(id),
    grade           VARCHAR(5)  NOT NULL,
    graded_by       BIGINT      REFERENCES faculty(id),
    graded_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
