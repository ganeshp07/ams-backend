package repositories

import "ams-backend/models"

type EmailSettingsRepository struct{}

func NewEmailSettingsRepository() *EmailSettingsRepository { return &EmailSettingsRepository{} }

func (r *EmailSettingsRepository) Get() (*models.EmailSettings, error) {
	var e models.EmailSettings
	err := DB.QueryRow(`
		SELECT id, smtp_host, smtp_port, smtp_username, smtp_password,
			from_email, from_name, is_enabled, created_at, updated_at
		FROM email_settings ORDER BY id LIMIT 1`,
	).Scan(&e.ID, &e.SmtpHost, &e.SmtpPort, &e.SmtpUsername, &e.SmtpPassword,
		&e.FromEmail, &e.FromName, &e.IsEnabled, &e.CreatedAt, &e.UpdatedAt)
	return &e, err
}

func (r *EmailSettingsRepository) Update(e *models.EmailSettings) error {
	// Only update password if a new one was provided (non-empty)
	if e.SmtpPassword != "" {
		_, err := DB.Exec(`
			UPDATE email_settings SET smtp_host=$1, smtp_port=$2, smtp_username=$3,
				smtp_password=$4, from_email=$5, from_name=$6, is_enabled=$7, updated_at=NOW()
			WHERE id=$8`,
			e.SmtpHost, e.SmtpPort, e.SmtpUsername, e.SmtpPassword,
			e.FromEmail, e.FromName, e.IsEnabled, e.ID)
		return err
	}
	_, err := DB.Exec(`
		UPDATE email_settings SET smtp_host=$1, smtp_port=$2, smtp_username=$3,
			from_email=$4, from_name=$5, is_enabled=$6, updated_at=NOW()
		WHERE id=$7`,
		e.SmtpHost, e.SmtpPort, e.SmtpUsername,
		e.FromEmail, e.FromName, e.IsEnabled, e.ID)
	return err
}
