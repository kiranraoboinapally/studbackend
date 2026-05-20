package swagger

import (
	"net/http"
	"os"
)

type Handler struct {
	specPath string
}

func NewHandler(specPath string) *Handler {
	return &Handler{specPath: specPath}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Serve the raw YAML spec
	mux.HandleFunc("/swagger/doc.yaml", h.serveYAML)
	// Serve Swagger UI
	mux.HandleFunc("/swagger/", h.serveUI)
	mux.HandleFunc("/swagger", h.serveUI)
}

func (h *Handler) serveYAML(w http.ResponseWriter, r *http.Request) {
	data, err := os.ReadFile(h.specPath)
	if err != nil {
		http.Error(w, "spec not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/yaml")
	w.Header().Set("Cache-Control", "no-cache")
	w.Write(data)
}

func (h *Handler) serveUI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(swaggerUIHTML))
}

const swaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>University ERP API - Swagger UI</title>
  <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  <style>
    html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
    *, *:before, *:after { box-sizing: inherit; }
    body { margin:0; background: #fafafa; }
    .topbar { display: none; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
  <script>
    window.onload = function() {
      const ui = SwaggerUIBundle({
        url: "/swagger/doc.yaml",
        dom_id: '#swagger-ui',
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout",
        persistAuthorization: true,
        displayRequestDuration: true,
        filter: true,
        tryItOutEnabled: true
      });
      window.ui = ui;
    };
  </script>
</body>
</html>
`
