package main

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/labstack/echo/v4"
)

type DNSRequest struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
}

func main() {
	e := echo.New()
	authToken := os.Getenv("AUTH_TOKEN")
	domainBase := os.Getenv("DOMAIN_BASE")
	if domainBase == "" {
		domainBase = "dns.local"
	}
	configPath := "/etc/dnsmasq.conf"

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

		changed, err := updateDNSConfig(req.IP, req.Hostname, configPath, domainBase)
		if err != nil {
			return c.String(http.StatusInternalServerError, "Failed to update config")
		}
		if !changed {
			return c.String(http.StatusOK, "No changes needed; config already up-to-date")
		}
		return c.String(http.StatusOK, fmt.Sprintf("DNS updated for %s.%s to %s", req.Hostname, domainBase, req.IP))
	})

	e.Logger.Fatal(e.Start(":8080"))
}

func updateDNSConfig(ip, hostname, configPath, domainBase string) (bool, error) {
	content, err := os.ReadFile(configPath)
	if err != nil {
		return false, err
	}

	lines := strings.Split(string(content), "\n")
	newEntry := fmt.Sprintf("address=/%s.%s/%s", hostname, domainBase, ip)
	updated := false
	changed := false

	for i, line := range lines {
		if strings.Contains(line, fmt.Sprintf("/%s.%s/", hostname, domainBase)) {
			if line != newEntry {
				lines[i] = newEntry
				changed = true
			}
			updated = true
			break
		}
	}

	if !updated {
		lines = append(lines, newEntry)
		changed = true
	}

	if !changed {
		return false, nil
	}

	err = os.WriteFile(configPath, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		return false, err
	}

	// Kill dnsmasq
	if err := exec.Command("pkill", "dnsmasq").Run(); err != nil {
		return false, fmt.Errorf("failed to kill dnsmasq: %w", err)
	}

	// Restart dnsmasq
	// Use full path if needed (e.g., /usr/sbin/dnsmasq) and required args if any
	cmd := exec.Command("dnsmasq", "--no-daemon")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	go cmd.Run() // Run in background so it doesn't block

	return true, nil
}
