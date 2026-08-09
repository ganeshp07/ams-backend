CREATE TABLE IF NOT EXISTS email_settings (
    id            BIGSERIAL PRIMARY KEY,
    smtp_host     VARCHAR(255) NOT NULL DEFAULT 'smtp.gmail.com',
    smtp_port     INT          NOT NULL DEFAULT 587,
    smtp_username VARCHAR(255) NOT NULL DEFAULT '',
    smtp_password VARCHAR(500) NOT NULL DEFAULT '',
    from_email    VARCHAR(255) NOT NULL DEFAULT '',
    from_name     VARCHAR(255) NOT NULL DEFAULT 'Academic Management System',
    is_enabled    BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
INSERT INTO email_settings DEFAULT VALUES ON CONFLICT DO NOTHING;
