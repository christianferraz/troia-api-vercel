package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/christianferraz/pkg/models"
	"github.com/skip2/go-qrcode"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var certificado models.CertificadoDTO
		json.NewDecoder(r.Body).Decode(&certificado)
		fmt.Fprintf(w, "Certificado recebido: %+v\n", certificado)
		pdfData, err := GerarCertificadoPDF(certificado)
		if err != nil {
			slog.Error("Erro ao gerar certificado PDF", "error", err)
			http.Error(w, "Erro ao gerar certificado PDF"+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"certificado_%s.pdf\"", strings.ReplaceAll(certificado.Nome, " ", "_")))
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(pdfData)))
		if _, err := w.Write(pdfData); err != nil {
			slog.Error("Erro ao escrever resposta PDF", "error", err)
			return
		}
		return
	}
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func GerarCertificadoPDF(participante models.CertificadoDTO) ([]byte, error) {
	// Lê o template SVG
	nomeArquivo := "certificado_palestrante.svg"

	svgTemplate, err := os.ReadFile(nomeArquivo)
	if err != nil {
		return nil, err
	}

	// Gera QR code
	qrCodeData, err := gerarQRCode(participante)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar QR code: %v", err)
	}

	// Substitui placeholders
	svgStr := string(svgTemplate)
	svgStr = strings.Replace(svgStr, "NOME_AQUI", participante.Nome, 1)
	if participante.Tipo != "ouvinte" {
		svgStr = strings.Replace(svgStr, "XXXXXXXX", participante.Tipo+" - 24h", 1)
	}

	// Adiciona o QR code no SVG
	svgStr = strings.Replace(svgStr, "QRCODE_AQUI", qrCodeData, 1)

	// Converte SVG para PDF usando rsvg-convert
	cmd := exec.Command("rsvg-convert", "-f", "pdf", "-")
	cmd.Stdin = bytes.NewBufferString(svgStr)

	var out bytes.Buffer
	cmd.Stdout = &out

	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

// gerarQRCode gera um QR code com informações do participante
func gerarQRCode(participante models.CertificadoDTO) (string, error) {

	qrContent := fmt.Sprintf("https://troia2025.sbop.org.br/verificar-certificado?id=%s", participante.ID)

	png, err := qrcode.Encode(qrContent, qrcode.Medium, 200)
	if err != nil {
		return "", err
	}

	// Converte para base64 para embed no SVG
	encoded := base64.StdEncoding.EncodeToString(png)

	// Retorna como data URI para usar no SVG
	return fmt.Sprintf("data:image/png;base64,%s", encoded), nil
}
