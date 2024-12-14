package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/sawatkins/upfast-tf/database"
)

// Test NotFound returns 404 and return html
func TestNotFound(t *testing.T) {
	engine := html.New("../templates", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Use(NotFound)

	req := httptest.NewRequest(http.MethodGet, "/non-existent-route", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	// Check if the status code is 404
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code 404, got %d", resp.StatusCode)
	}

	// Check that the content type response is html
	if contentType := resp.Header.Get("Content-Type"); contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected content type 'text/html; charset=utf-8', got '%s'", contentType)
	}
}

// Test Index route existance, status code, and content-type
func TestIndex(t *testing.T) {
	database.InitDB(":memory:")
	database.InitPlayerSessionTable()

	engine := html.New("../templates", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Get("/", Index)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	// Check if the status code is 200
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	// Check that the content type response is html
	if contentType := resp.Header.Get("Content-Type"); contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected content type 'text/html; charset=utf-8', got '%s'", contentType)
	}
}

// Test About route existance, status code, and content-type
func TestAbout(t *testing.T) {
	engine := html.New("../templates", ".html")
	app := fiber.New(fiber.Config{Views: engine})
	app.Get("/about", About)

	req := httptest.NewRequest(http.MethodGet, "/about", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to send request: %v", err)
	}

	// Check if the status code is 200
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}

	// Check that the content type response is html
	if contentType := resp.Header.Get("Content-Type"); contentType != "text/html; charset=utf-8" {
		t.Errorf("Expected content type 'text/html; charset=utf-8', got '%s'", contentType)
	}
}

func TestGetServerIPs(t *testing.T) {
	type args struct {
		c *fiber.Ctx
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GetServerIPs(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("GetServerIPs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
