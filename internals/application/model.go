package application

import (
	"time"
)

type Application struct {
	ApplicationID   string     `bson:"application_id" json:"application_id"`
	TeamName        string     `bson:"team_name" json:"team_name"`
	PMName          string     `bson:"pm_name" json:"pm_name"`
	PMEmail         string     `bson:"pm_email" json:"pm_email"`
	PMContact       string     `bson:"pm_contact" json:"pm_contact"`
	AlternateNumber string     `bson:"alternate_number" json:"alternate_number"`
	Domain          string     `bson:"domain" json:"domain"`
	Teammates       []Teammate `bson:"teammates" json:"teammates"`
	ProposalPDFURL  string     `bson:"proposal_pdf_url" json:"proposal_pdf_url"`
	Status          string     `bson:"status" json:"status"`
	CreatedAt       time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `bson:"updated_at" json:"updated_at"`
}

type Teammate struct {
	Name      string `bson:"name" json:"name"`
	Email     string `bson:"email" json:"email"`
	Role      string `bson:"role" json:"role"`
	ResumeURL string `bson:"resume_url" json:"resume_url"`
	Portfolio string `bson:"portfolio,omitempty" json:"portfolio,omitempty"`
	Github    string `bson:"github,omitempty" json:"github,omitempty"`
}
