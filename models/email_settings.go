package models

import "time"

type EmailSettings struct {
	ID           int64     `json:"id"`
	SmtpHost     string    `json:"smtp_host"`
	SmtpPort     int       `json:"smtp_port"`
	SmtpUsername string    `json:"smtp_username"`
	SmtpPassword string    `json:"-"`          // never serialised
	FromEmail    string    `json:"from_email"`
	FromName     string    `json:"from_name"`
	IsEnabled    bool      `json:"is_enabled"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// EmailSettingsResponse is the safe version returned by the API (password masked).
type EmailSettingsResponse struct {
	ID             int64     `json:"id"`
	SmtpHost       string    `json:"smtp_host"`
	SmtpPort       int       `json:"smtp_port"`
	SmtpUsername   string    `json:"smtp_username"`
	SmtpPassword   string    `json:"smtp_password"` // masked
	FromEmail      string    `json:"from_email"`
	FromName       string    `json:"from_name"`
	IsEnabled      bool      `json:"is_enabled"`
	UpdatedAt      time.Time `json:"updated_at"`
}
