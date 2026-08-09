package services

import (
	"ams-backend/repositories"
	"fmt"
	"log"
	"net/smtp"
	"strconv"
)

type EmailService struct {
	emailRepo *repositories.EmailSettingsRepository
}

func NewEmailService(r *repositories.EmailSettingsRepository) *EmailService {
	return &EmailService{emailRepo: r}
}

func (s *EmailService) SendGradeNotification(to, studentName, courseName, courseCode, semesterName, grade string) {
	go func() {
		if err := s.sendGrade(to, studentName, courseName, courseCode, semesterName, grade); err != nil {
			log.Printf("grade email error (to %s): %v", to, err)
		}
	}()
}

func (s *EmailService) sendGrade(to, studentName, courseName, courseCode, semesterName, grade string) error {
	settings, err := s.emailRepo.Get()
	if err != nil {
		return fmt.Errorf("load email settings: %w", err)
	}

	if !settings.IsEnabled {
		log.Printf("[EMAIL DISABLED] Grade notification for %s: %s in %s (%s) = %s",
			studentName, grade, courseName, semesterName, to)
		return nil
	}

	subject := "Grade Updated - " + courseName
	body := fmt.Sprintf(`Dear %s,

Your grade for %s has been updated.

Course: %s
Course Code: %s
Semester: %s
Grade: %s

Please login to the Academic Management System to view your academic record.

Regards,
%s`, studentName, courseName, courseName, courseCode, semesterName, grade, settings.FromName)

	msg := fmt.Sprintf("From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		settings.FromName, settings.FromEmail, to, subject, body)

	auth := smtp.PlainAuth("", settings.SmtpUsername, settings.SmtpPassword,
		settings.SmtpHost)
	addr := settings.SmtpHost + ":" + strconv.Itoa(settings.SmtpPort)
	return smtp.SendMail(addr, auth, settings.FromEmail, []string{to}, []byte(msg))
}

func (s *EmailService) TestConnection(testTo string) error {
	settings, err := s.emailRepo.Get()
	if err != nil {
		return fmt.Errorf("load settings: %w", err)
	}
	if settings.SmtpHost == "" || settings.SmtpUsername == "" {
		return fmt.Errorf("SMTP settings are not configured")
	}

	msg := fmt.Sprintf("From: %s <%s>\r\nTo: %s\r\nSubject: Test Email - AMS\r\n\r\nThis is a test email from the Academic Management System.",
		settings.FromName, settings.FromEmail, testTo)

	auth := smtp.PlainAuth("", settings.SmtpUsername, settings.SmtpPassword, settings.SmtpHost)
	addr := settings.SmtpHost + ":" + strconv.Itoa(settings.SmtpPort)
	return smtp.SendMail(addr, auth, settings.FromEmail, []string{testTo}, []byte(msg))
}
