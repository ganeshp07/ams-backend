package services

import (
	"ams-backend/repositories"
	"crypto/tls"
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
		log.Printf("[EMAIL DISABLED] Grade notification for %s: %s in %s = %s",
			studentName, grade, courseName, to)
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

	return s.sendMail(settings.SmtpHost, settings.SmtpPort, settings.SmtpUsername,
		settings.SmtpPassword, settings.FromEmail, settings.FromName, to, subject, body)
}

func (s *EmailService) TestConnection(testTo string) error {
	settings, err := s.emailRepo.Get()
	if err != nil {
		return fmt.Errorf("load settings: %w", err)
	}
	if settings.SmtpHost == "" || settings.SmtpUsername == "" {
		return fmt.Errorf("SMTP settings are not configured")
	}
	subject := "Test Email - AMS"
	body := "This is a test email from the Academic Management System."
	return s.sendMail(settings.SmtpHost, settings.SmtpPort, settings.SmtpUsername,
		settings.SmtpPassword, settings.FromEmail, settings.FromName, testTo, subject, body)
}

// sendMail tries SSL port 465 first (works on Render), then STARTTLS port 587 (works locally).
func (s *EmailService) sendMail(host string, port int, username, password, fromEmail, fromName, to, subject, body string) error {
	msg := fmt.Sprintf("From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\n\r\n%s",
		fromName, fromEmail, to, subject, body)

	// Try SSL on port 465 first — Render allows this
	if err := s.sendSSL(host, username, password, fromEmail, to, []byte(msg)); err == nil {
		return nil
	}

	// Fallback: STARTTLS on port 587 — works locally and most other hosts
	addr := host + ":" + strconv.Itoa(port)
	auth := smtp.PlainAuth("", username, password, host)
	return smtp.SendMail(addr, auth, fromEmail, []string{to}, []byte(msg))
}

// sendSSL connects using TLS on port 465 — needed on Render free tier.
func (s *EmailService) sendSSL(host, username, password, from, to string, msg []byte) error {
	tlsConfig := &tls.Config{
		ServerName: host,
	}
	conn, err := tls.Dial("tcp", host+":465", tlsConfig)
	if err != nil {
		return fmt.Errorf("SSL dial: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("SMTP client: %w", err)
	}
	defer client.Quit()

	if err := client.Auth(smtp.PlainAuth("", username, password, host)); err != nil {
		return fmt.Errorf("auth: %w", err)
	}
	if err := client.Mail(from); err != nil {
		return err
	}
	if err := client.Rcpt(to); err != nil {
		return err
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err := w.Write(msg); err != nil {
		return err
	}
	return w.Close()
}
