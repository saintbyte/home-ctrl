package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/saintbyte/home-ctrl/internal/config"
)

// CertificateHandler handles the root certificate download endpoint
type CertificateHandler struct {
	config *config.Config
}

// NewCertificateHandler creates a new certificate handler
func NewCertificateHandler(cfg *config.Config) *CertificateHandler {
	return &CertificateHandler{config: cfg}
}

// SetupRoutes sets up certificate-related routes
func (h *CertificateHandler) SetupRoutes(router *gin.RouterGroup) {
	router.GET("/cert", h.downloadRootCert)
}

// downloadRootCert handles GET /cert - downloads the root certificate
func (h *CertificateHandler) downloadRootCert(c *gin.Context) {
	certPath := h.config.Server.CACert
	if certPath == "" {
		c.JSON(http.StatusNotFound, gin.H{
			"error":  "root certificate is not configured",
			"detail": "set ca_cert in the server section of the config",
		})
		return
	}
	// Serve the certificate file as a download attachment
	c.FileAttachment(certPath, "home-ctrl-root-ca.crt")
}