package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/christianferraz/api/models"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var certificado models.CertificadoDTO
		json.NewDecoder(r.Body).Decode(&certificado)
		fmt.Fprintf(w, "Certificado recebido: %+v\n", certificado)
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}
