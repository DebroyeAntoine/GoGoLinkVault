package scraper

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestScrapeMetadata(t *testing.T) {
	// Serveur de test avec du HTML contenant les meta tags
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		html := `
			<html>
			<head>
				<title>Test Page</title>
				<meta name="description" content="This is a test description.">
				<meta property="og:image" content="https://example.com/image.jpg">
			</head>
			<body><p>Hello World!</p></body>
			</html>`
		w.Write([]byte(html))
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	metadata, err := FetchMetadata(server.URL)

	assert.NoError(t, err)
	assert.Equal(t, "Test Page", metadata.Title)
	assert.Equal(t, "This is a test description.", metadata.Description)
	assert.Equal(t, "https://example.com/image.jpg", metadata.Image)
}
