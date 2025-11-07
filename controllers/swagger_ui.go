package controllers

import (
	"io/ioutil"
	"net/http"

	"github.com/beego/beego/v2/server/web"
)

type SwaggerUIController struct {
	web.Controller
}

// UI serves a minimal Swagger UI that loads /swagger/doc.json
func (s *SwaggerUIController) UI() {
	html := `<!DOCTYPE html>
<html>
  <head>
    <meta charset="UTF-8">
    <title>API Docs</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@4/swagger-ui.css" />
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@4/swagger-ui-bundle.js"></script>
    <script>
      window.onload = function() {
        const ui = SwaggerUIBundle({
          url: '/swagger/doc.json',
          dom_id: '#swagger-ui',
          presets: [SwaggerUIBundle.presets.apis],
        })
      }
    </script>
  </body>
</html>`

	s.Ctx.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
	s.Ctx.ResponseWriter.Write([]byte(html))
}

// Doc serves the generated swagger.json file
func (s *SwaggerUIController) Doc() {
	data, err := ioutil.ReadFile("docs/swagger.json")
	if err != nil {
		s.Ctx.ResponseWriter.WriteHeader(http.StatusNotFound)
		s.Ctx.ResponseWriter.Write([]byte("swagger doc not found"))
		return
	}
	s.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
	s.Ctx.ResponseWriter.Write(data)
}
