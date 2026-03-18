package application

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"time"

	"github.com/suhas-developer07/EdwinNova-Server/internals/infrastructure/mail"
	"github.com/xuri/excelize/v2"
)

type Service interface {
	CreateApplication(ctx context.Context, app *Application) error
	ExportApplications(ctx context.Context) ([]byte, error)
	ExportApplicationsCSV(ctx context.Context) ([]byte, error)
	GetAllApplications(ctx context.Context) ([]Application, error)
}

type service struct {
	repo Repository
	// publisher EmailPublisher
	smtp *mail.ResendClient
}

func NewService(repo Repository, smtp *mail.ResendClient) Service {
	return &service{
		repo: repo,
		smtp: smtp,
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

	log.Printf("Application created for team %s with PM %s", app.TeamName, app.PMEmail)

	emailBody, err := mail.BuildRegistrationEmailBody(app.TeamName, app.PMName, app.PMEmail, app.PMContact, app.ApplicationID, app.CreatedAt)
	if err != nil {
		log.Printf("Failed to build registration email body for team %s: %v", app.TeamName, err)
		return err
	}
	err = s.smtp.Send(app.PMEmail, "Your Hackothon registration is successfull", emailBody)
	if err != nil {
		log.Printf("Failed to send registration email to %s: %v", app.PMEmail, err)
		return err
	}
	log.Printf("Sent registration email to %s for team %s", app.PMEmail, app.TeamName)
	return nil
}
func (s *service) ExportApplications(ctx context.Context) ([]byte, error) {

	apps, err := s.repo.GetAllApplications(ctx)
	if err != nil {
		return nil, err
	}

	file := excelize.NewFile()
	sheet := "Applications"
	file.SetSheetName("Sheet1", sheet)

	headers := []string{
		"ApplicationID",
		"TeamName",
		"PMName",
		"PMEmail",
		"PMContact",
		"AlternateNumber",
		"Domain",
		"Status",
		"ProposalPDFURL",
		"TeammateName",
		"TeammateEmail",
		"Role",
	}

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		file.SetCellValue(sheet, cell, h)
	}

	row := 2

	for _, app := range apps {

		if len(app.Teammates) == 0 {

			file.SetCellValue(sheet, fmt.Sprintf("A%d", row), app.ApplicationID)
			file.SetCellValue(sheet, fmt.Sprintf("B%d", row), app.TeamName)
			file.SetCellValue(sheet, fmt.Sprintf("C%d", row), app.PMName)
			file.SetCellValue(sheet, fmt.Sprintf("D%d", row), app.PMEmail)
			file.SetCellValue(sheet, fmt.Sprintf("E%d", row), app.PMContact)
			file.SetCellValue(sheet, fmt.Sprintf("F%d", row), app.AlternateNumber)
			file.SetCellValue(sheet, fmt.Sprintf("G%d", row), app.Domain)
			file.SetCellValue(sheet, fmt.Sprintf("H%d", row), app.Status)
			file.SetCellValue(sheet, fmt.Sprintf("I%d", row), app.ProposalPDFURL)

			row++
			continue
		}

		for _, member := range app.Teammates {

			file.SetCellValue(sheet, fmt.Sprintf("A%d", row), app.ApplicationID)
			file.SetCellValue(sheet, fmt.Sprintf("B%d", row), app.TeamName)
			file.SetCellValue(sheet, fmt.Sprintf("C%d", row), app.PMName)
			file.SetCellValue(sheet, fmt.Sprintf("D%d", row), app.PMEmail)
			file.SetCellValue(sheet, fmt.Sprintf("E%d", row), app.PMContact)
			file.SetCellValue(sheet, fmt.Sprintf("F%d", row), app.AlternateNumber)
			file.SetCellValue(sheet, fmt.Sprintf("G%d", row), app.Domain)
			file.SetCellValue(sheet, fmt.Sprintf("H%d", row), app.Status)
			file.SetCellValue(sheet, fmt.Sprintf("I%d", row), app.ProposalPDFURL)

			file.SetCellValue(sheet, fmt.Sprintf("J%d", row), member.Name)
			file.SetCellValue(sheet, fmt.Sprintf("K%d", row), member.Email)
			file.SetCellValue(sheet, fmt.Sprintf("L%d", row), member.Role)

			row++
		}
	}

	buf := new(bytes.Buffer)

	if err := file.Write(buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *service) ExportApplicationsCSV(ctx context.Context) ([]byte, error) {

	apps, err := s.repo.GetAllApplications(ctx)
	if err != nil {
		return nil, err
	}

	buffer := &bytes.Buffer{}
	writer := csv.NewWriter(buffer)

	header := []string{
		"ApplicationID",
		"TeamName",
		"PMName",
		"PMEmail",
		"PMContact",
		"AlternateNumber",
		"Domain",
		"Status",
		"ProposalPDFURL",
		"TeammateName",
		"TeammateEmail",
		"Role",
	}

	writer.Write(header)

	for _, app := range apps {

		if len(app.Teammates) == 0 {

			row := []string{
				app.ApplicationID,
				app.TeamName,
				app.PMName,
				app.PMEmail,
				app.PMContact,
				app.AlternateNumber,
				app.Domain,
				app.Status,
				app.ProposalPDFURL,
				"",
				"",
				"",
				"",
				"",
				"",
			}

			writer.Write(row)
			continue
		}

		for _, member := range app.Teammates {

			row := []string{
				app.ApplicationID,
				app.TeamName,
				app.PMName,
				app.PMEmail,
				app.PMContact,
				app.AlternateNumber,
				app.Domain,
				app.Status,
				app.ProposalPDFURL,
				member.Name,
				member.Email,
				member.Role,
			}

			writer.Write(row)
		}
	}

	writer.Flush()

	return buffer.Bytes(), nil
}

func (s *service) GetAllApplications(ctx context.Context) ([]Application, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	return s.repo.GetAllApplications(ctx)
}