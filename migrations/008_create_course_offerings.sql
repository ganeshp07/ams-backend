CREATE TABLE IF NOT EXISTS course_offerings (
    id            BIGSERIAL PRIMARY KEY,
    course_id     BIGINT       NOT NULL REFERENCES courses(id),
    semester_id   BIGINT       NOT NULL REFERENCES semesters(id),
    faculty_id    BIGINT       NOT NULL REFERENCES faculty(id),
    academic_year VARCHAR(20)  NOT NULL,
    section       VARCHAR(10)  NOT NULL DEFAULT 'A',
    is_active     BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE(course_id, semester_id, section)
);
