package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	service   Service
	uploadDir string
}

func NewHandler(service Service, uploadDir string) *Handler {
	return &Handler{
		service:   service,
		uploadDir: uploadDir,
	}
}

type createApplicationRequest struct {
	TeamName        string            `json:"team_name"`
	PMName          string            `json:"pm_name"`
	PMEmail         string            `json:"pm_email"`
	PMContact       string            `json:"pm_contact"`
	AlternateNumber string            `json:"alternate_number"`
	Domain          string            `json:"domain"`
	Teammates       []teammatePayload `json:"teammates"`
}

type teammatePayload struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Github string `json:"github"`
}

func (h *Handler) CreateApplication(c echo.Context) error {
	ctx := c.Request().Context()

	var req createApplicationRequest
	if err := json.Unmarshal([]byte(c.FormValue("payload")), &req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid payload json")
	}

	proposalFile, err := c.FormFile("proposal")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "proposal pdf is required")
	}

	form, err := c.MultipartForm()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid multipart form")
	}
	resumeFiles := form.File["resumes"]

	if err := validateCreateApplicationRequest(&req, proposalFile, resumeFiles); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := os.MkdirAll(h.uploadDir, 0o755); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "could not prepare upload directory")
	}

	proposalURL, err := h.savePDF(proposalFile, "proposals")
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to store proposal pdf")
	}

	var resumeURLs []string
	for _, f := range resumeFiles {
		url, err := h.savePDF(f, "resumes")
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to store resume pdf")
		}
		resumeURLs = append(resumeURLs, url)
	}

	if len(resumeURLs) != len(req.Teammates) {
		return echo.NewHTTPError(http.StatusBadRequest, "number of resumes must match number of teammates")
	}

	teammates := make([]Teammate, len(req.Teammates))
	for i, t := range req.Teammates {
		teammates[i] = Teammate{
			Name:      t.Name,
			Email:     t.Email,
			Role:      t.Role,
			ResumeURL: resumeURLs[i],
			Github:    t.Github,
		}
	}

	app := &Application{
		ApplicationID:   primitive.NewObjectID().Hex(),
		TeamName:        req.TeamName,
		PMName:          req.PMName,
		PMEmail:         req.PMEmail,
		PMContact:       req.PMContact,
		AlternateNumber: req.AlternateNumber,
		Domain:          req.Domain,
		Teammates:       teammates,
		ProposalPDFURL:  proposalURL,
		Status:          "pending",
		CreatedAt:       time.Time{},
		UpdatedAt:       time.Time{},
	}

	if err := h.service.CreateApplication(ctx, app); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create application")
	}

	return c.JSON(http.StatusCreated, app)
}

func validateCreateApplicationRequest(req *createApplicationRequest, proposal *multipart.FileHeader, resumes []*multipart.FileHeader) error {
	if req.TeamName == "" ||
		req.PMName == "" ||
		req.PMEmail == "" ||
		req.PMContact == "" ||
		req.Domain == "" {
		return errors.New("missing required fields")
	}

	if !isValidEmail(req.PMEmail) {
		return errors.New("invalid pm_email")
	}

	if len(req.Teammates) == 0 {
		return errors.New("at least one teammate is required")
	}

	for _, t := range req.Teammates {
		if t.Name == "" || t.Email == "" || t.Role == "" {
			return errors.New("each teammate must have name, email and role")
		}
		if !isValidEmail(t.Email) {
			return errors.New("invalid teammate email")
		}
	}

	if err := validatePDFHeader(proposal); err != nil {
		return err
	}
	for _, r := range resumes {
		if err := validatePDFHeader(r); err != nil {
			return err
		}
	}

	return nil
}

func isValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func (h *Handler) savePDF(fileHeader *multipart.FileHeader, subdir string) (string, error) {
	filename := filepath.Base(fileHeader.Filename)
	ext := filepath.Ext(filename)
	if ext != ".pdf" {
		return "", errors.New("only pdf files are allowed")
	}

	targetDir := filepath.Join(h.uploadDir, subdir)
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return "", err
	}

	newName := fmt.Sprintf("%s_%s", primitive.NewObjectID().Hex(), filename)
	targetPath := filepath.Join(targetDir, newName)

	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	dst, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return fmt.Sprintf("/uploads/%s/%s", subdir, newName), nil
}

func validatePDFHeader(fileHeader *multipart.FileHeader) error {
	filename := fileHeader.Filename
	if filepath.Ext(filename) != ".pdf" {
		return errors.New("file must be a pdf")
	}
	// Size checks can be added here if max upload size is enforced.
	return nil
}
