CREATE TABLE IF NOT EXISTS semesters (
    id                 BIGSERIAL PRIMARY KEY,
    academic_year      VARCHAR(20)  NOT NULL,
    semester_number    INT          NOT NULL,
    name               VARCHAR(100) NOT NULL,
    start_date         DATE,
    end_date           DATE,
    registration_start TIMESTAMPTZ,
    registration_end   TIMESTAMPTZ,
    is_active          BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
