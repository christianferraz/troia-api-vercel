package models

import "github.com/google/uuid"

type CertificadoDTO struct {
	ID   uuid.UUID `json:"id,omitempty"`
	Nome string    `json:"nome"`
	Tipo string    `json:"tipo"`
}

type Certificado struct {
	UserID uuid.UUID `json:"user_id"`
	Tipo   string    `json:"tipo"`
}
