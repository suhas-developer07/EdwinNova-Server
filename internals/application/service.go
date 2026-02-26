package application

import (
	"context"
	"log"
	"time"

	"github.com/suhas-developer07/EdwinNova-Server/internals/email"
)

type Service interface {
	CreateApplication(ctx context.Context, app *Application) error
}

type EmailPublisher interface {
	Publish(ctx context.Context, job email.EmailJob) error
}

type service struct {
	repo      Repository
	publisher EmailPublisher
}

func NewService(repo Repository,publisher EmailPublisher) Service {
	return &service{
		repo: repo,
		publisher: publisher,
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

	err = s.publisher.Publish(ctx, email.EmailJob{
		To:       app.PMEmail,
		Subject:  "Hackathon Registration Successfull",
		Template: "application_created",
		Data: map[string]interface{}{
			"name": app.TeamName,
		},
	})

	if err != nil {
		log.Printf("Appication successfull email failed,%v",app.TeamName)
	}
	return err
}
