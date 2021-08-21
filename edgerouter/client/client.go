package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

type Client struct {
	baseUrl string
	client  *http.Client
}

// NewClient creates a new client with the given TLS configuration, most useful for setting InsecureSkipVerify or RootCAs
func NewClient(tlsConfig *tls.Config, baseUrl, username, password string) (*Client, error) {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("error creating cookie jar: %w", err)
	}
	c := &Client{
		baseUrl: baseUrl,
		client: &http.Client{
			CheckRedirect: disableFollowingRedirects,
			Jar:           cookieJar,
			Transport:     &http.Transport{TLSClientConfig: tlsConfig},
		},
	}
	if err := c.login(username, password); err != nil {
		return nil, fmt.Errorf("error logging in: %w", err)
	}

	return c, nil
}

// disableFollowingRedirects causes the HTTP client to return the original response instead of e.g. following the redirect after login
func disableFollowingRedirects(*http.Request, []*http.Request) error {
	return http.ErrUseLastResponse
}

func (c *Client) login(username, password string) error {
	loginForm := url.Values{}
	loginForm.Set("username", username)
	loginForm.Set("password", password)

	resp, err := c.client.PostForm(c.baseUrl, loginForm)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 303 {
		return fmt.Errorf("unexpected response code %d: %+v", resp.StatusCode, resp)
	}
	if _, err := io.ReadAll(resp.Body); err != nil {
		return fmt.Errorf("error reading login response: %w", err)
	}
	log.Printf("[TRACE] Login suceeded as %s", username)

	return nil
}

func (c *Client) Post(ctx context.Context, path string, input, output interface{}) error {
	body, err := json.Marshal(input)
	log.Printf("[DEBUG] Sending request: %s", body)
	if err != nil {
		return fmt.Errorf("error marshalling input body: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseUrl+path, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("could not create HTTP request: %w", err)
	}
	// X-Requested-With: XMLHTTPRequest is required to avoid a 403
	req.Header.Add("X-Requested-With", "XMLHTTPRequest")

	resp, err := c.client.Do(req)
	if err != nil {
		log.Printf("[WARN] Got unknwon error: %s", err)
		return fmt.Errorf("could not perform HTTP request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("[DEBUG] Got response code %d", resp.StatusCode)
		return fmt.Errorf("unexpected response code %d: %+v", resp.StatusCode, resp)
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	bodyString := buf.String()
	log.Printf("[DEBUG] Got response: %s", bodyString)
	if err := json.NewDecoder(buf).Decode(&output); err != nil {
		return fmt.Errorf("error unmarshalling response body: %w", err)
	}

	return nil
}
