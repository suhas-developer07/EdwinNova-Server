package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"regexp"
	"time"

	storage "github.com/suhas-developer07/EdwinNova-Server/internals/infrastructure/s3"

	"github.com/labstack/echo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	service Service
	storage *storage.S3Storage
}

func NewHandler(service Service, storage *storage.S3Storage) *Handler {
	return &Handler{
		service: service,
		storage: storage,
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
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Github    string `json:"github"`
	Portfolio string `json:"portfolio"`
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

	applicationID := primitive.NewObjectID().Hex()

	/* Upload proposal */
	proposalKey := fmt.Sprintf("applications/%s/proposal.pdf", applicationID)

	proposalURL, err := h.storage.UploadFile(ctx, proposalFile, proposalKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to upload proposal")
	}

	/* Upload resumes */
	var resumeURLs []string

	for i, f := range resumeFiles {

		key := fmt.Sprintf(
			"applications/%s/resume_%d.pdf",
			applicationID,
			i+1,
		)

		url, err := h.storage.UploadFile(ctx, f, key)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to upload resume")
		}

		resumeURLs = append(resumeURLs, url)
	}

	if len(resumeURLs) != len(req.Teammates) {
		return echo.NewHTTPError(http.StatusBadRequest, "resumes must match teammates")
	}

	teammates := make([]Teammate, len(req.Teammates))

	for i, t := range req.Teammates {

		teammates[i] = Teammate{
			Name:      t.Name,
			Email:     t.Email,
			Role:      t.Role,
			ResumeURL: resumeURLs[i],
			Github:    t.Github,
			Portfolio: t.Portfolio,
		}
	}

	app := &Application{
		ApplicationID:   applicationID,
		TeamName:        req.TeamName,
		PMName:          req.PMName,
		PMEmail:         req.PMEmail,
		PMContact:       req.PMContact,
		AlternateNumber: req.AlternateNumber,
		Domain:          req.Domain,
		Teammates:       teammates,
		ProposalPDFURL:  proposalURL,
		Status:          "pending",
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := h.service.CreateApplication(ctx, app); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to create application")
	}

	return c.JSON(http.StatusCreated, app)
}

func validateCreateApplicationRequest(
	req *createApplicationRequest,
	proposal *multipart.FileHeader,
	resumes []*multipart.FileHeader,
) error {

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
		return errors.New("at least one teammate required")
	}

	for _, t := range req.Teammates {

		if t.Name == "" || t.Email == "" || t.Role == "" {
			return errors.New("each teammate must have name email role")
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

func validatePDFHeader(fileHeader *multipart.FileHeader) error {

	if filepath.Ext(fileHeader.Filename) != ".pdf" {
		return errors.New("file must be pdf")
	}

	return nil
}

func isValidEmail(email string) bool {

	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

	return re.MatchString(email)
}