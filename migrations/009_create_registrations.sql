CREATE TABLE IF NOT EXISTS registrations (
    id                 BIGSERIAL PRIMARY KEY,
    student_id         BIGINT      NOT NULL REFERENCES students(id),
    course_offering_id BIGINT      NOT NULL REFERENCES course_offerings(id),
    registered_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(student_id, course_offering_id)
);
