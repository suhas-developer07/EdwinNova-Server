package application

import (
	"context"
	"log"
	"time"

	"github.com/suhas-developer07/EdwinNova-Server/internals/email"
	"github.com/suhas-developer07/EdwinNova-Server/internals/infrastructure/mail"
)

type Service interface {
	CreateApplication(ctx context.Context, app *Application) error
}

// type EmailPublisher interface {
// 	Publish(ctx context.Context, job email.EmailJob) error
// }

type service struct {
	repo      Repository
	// publisher EmailPublisher
	smtp      *mail.SMTPClient
}

func NewService(repo Repository, smtp *mail.SMTPClient) Service {
	return &service{
		repo:      repo,
		smtp:      smtp,
	}
}

func (s *service) CreateApplication(ctx context.Context, app *Application) error {
	now := time.Now().UTC()
	if app.CreatedAt.IsZero() {
		app.CreatedAt = now
	}
	app.UpdatedAt = now
	err := s.repo.Create(ctx, app)
	if err != nil {
		return err
	}

	emailBody,err := email.BuildRegistrationEmailBody(app.TeamName,app.PMName,app.PMEmail,app.PMContact,app.ApplicationID,app.CreatedAt)
	if err != nil {
		log.Printf("Failed to build registration email body for team %s: %v", app.TeamName, err)
		return err
	}
	err = s.smtp.Send(app.PMEmail,"Your Hackothon registration is successfull",emailBody)
	if err != nil {
		log.Printf("Failed to send registration email to %s: %v", app.PMEmail, err)
		return err
	}
	return nil
}
