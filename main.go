package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
	"strings"
)

type DNSRequest struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
}

func main() {
	e := echo.New()
	authToken := os.Getenv("AUTH_TOKEN")
	configPath := "/app/dnsmasq.conf"

	e.POST("/update-dns", func(c echo.Context) error {
		if c.Request().Header.Get("Authorization") != authToken {
			return c.String(http.StatusUnauthorized, "Unauthorized")
		}

		var req DNSRequest
		if err := c.Bind(&req); err != nil {
			return c.String(http.StatusBadRequest, "Invalid request body")
		}

		if req.IP == "" || req.Hostname == "" {
			return c.String(http.StatusBadRequest, "IP and hostname required")
		}

		if err := updateDNSConfig(req.IP, req.Hostname, configPath); err != nil {
			return c.String(http.StatusInternalServerError, "Failed to update config")
		}

		return c.String(http.StatusOK, fmt.Sprintf("DNS updated for %s.nsa.local to %s", req.Hostname, req.IP))
	})

	e.Logger.Fatal(e.Start(":8080"))
}

func updateDNSConfig(ip, hostname, configPath string) error {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	newEntry := fmt.Sprintf("address=/%s.nsa.local/%s", hostname, ip)
	updated := false

	for i, line := range lines {
		if strings.Contains(line, fmt.Sprintf("/%s.nsa.local/", hostname)) {
			lines[i] = newEntry
			updated = true
			break
		}
	}

	if !updated {
		lines = append(lines, newEntry)
	}

	return os.WriteFile(configPath, []byte(strings.Join(lines, "\n")), 0644)
}
