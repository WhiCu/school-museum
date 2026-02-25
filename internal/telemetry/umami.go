package telemetry

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

// UmamiInfo содержит конфигурацию Umami, которая передаётся на фронтенд
// для подключения клиентского трекинг-скрипта.
type UmamiInfo struct {
	URL       string `json:"url"`
	WebsiteID string `json:"website_id"`
}

// UmamiOpts содержит параметры для резолва сайта в Umami.
type UmamiOpts struct {
	URL       string
	WebsiteID string
	Username  string
	Password  string
	Domain    string
}

// ResolveUmamiWebsite проверяет наличие сайта в Umami и при необходимости создаёт его.
// Возвращает UmamiInfo с URL и WebsiteID для подключения трекинг-скрипта на фронтенде.
func ResolveUmamiWebsite(ctx context.Context, opts UmamiOpts, log *slog.Logger) (*UmamiInfo, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	websiteID := opts.WebsiteID
	if websiteID == "" {
		id, err := ensureWebsite(ctx, client, opts, log)
		if err != nil {
			return nil, fmt.Errorf("umami: auto-create website: %w", err)
		}
		websiteID = id
	}

	return &UmamiInfo{
		URL:       opts.URL,
		WebsiteID: websiteID,
	}, nil
}

// ensureWebsite логинится в Umami, ищет сайт по домену; если не найден — создаёт.
// Повторяет попытки, т.к. Umami может быть ещё не готов при старте.
func ensureWebsite(ctx context.Context, client *http.Client, opts UmamiOpts, log *slog.Logger) (string, error) {
	const maxRetries = 10
	var lastErr error

	for i := range maxRetries {
		if i > 0 {
			log.Info("umami: waiting for service...", slog.Int("attempt", i+1))
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(3 * time.Second):
			}
		}

		token, err := login(ctx, client, opts.URL, opts.Username, opts.Password)
		if err != nil {
			lastErr = fmt.Errorf("login: %w", err)
			continue
		}

		id, err := findWebsite(ctx, client, opts.URL, token, opts.Domain)
		if err != nil {
			lastErr = fmt.Errorf("find website: %w", err)
			continue
		}
		if id != "" {
			log.Info("umami: found existing website", slog.String("id", id), slog.String("domain", opts.Domain))
			return id, nil
		}

		id, err = createWebsite(ctx, client, opts.URL, token, opts.Domain)
		if err != nil {
			lastErr = fmt.Errorf("create website: %w", err)
			continue
		}
		log.Info("umami: created website", slog.String("id", id), slog.String("domain", opts.Domain))
		return id, nil
	}

	return "", fmt.Errorf("after %d attempts: %w", maxRetries, lastErr)
}

func login(ctx context.Context, client *http.Client, baseURL, username, password string) (string, error) {
	body, _ := json.Marshal(map[string]string{"username": username, "password": password})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/api/auth/login", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("status %d: %s", resp.StatusCode, string(data))
	}

	var result struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.Token, nil
}

func findWebsite(ctx context.Context, client *http.Client, baseURL, token, domain string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"/api/websites", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("status %d: %s", resp.StatusCode, string(data))
	}

	var result struct {
		Data []struct {
			ID     string `json:"id"`
			Domain string `json:"domain"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	for _, w := range result.Data {
		if w.Domain == domain {
			return w.ID, nil
		}
	}
	return "", nil
}

func createWebsite(ctx context.Context, client *http.Client, baseURL, token, domain string) (string, error) {
	body, _ := json.Marshal(map[string]string{"domain": domain, "name": domain})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/api/websites", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("status %d: %s", resp.StatusCode, string(data))
	}

	var result struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.ID, nil
}
