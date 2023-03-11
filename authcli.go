package authcli

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"strings"
)

type Client struct {
	Provider    Provider
	Credentials Credentials
	RedirectURL string
	Scopes      []string
}

type Credentials struct {
	ClientID     string
	ClientSecret string
}

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
}

type GenericError struct {
	Error     string `json:"error"`
	ErrorDesc string `json:"error_description"`
}

func New(provider Provider, clientID, clientSecret, redirectURL string, opts ...Option) *Client {
	c := Client{
		Provider: provider,
		Credentials: Credentials{
			ClientID:     clientID,
			ClientSecret: clientSecret,
		},
		RedirectURL: redirectURL,
	}
	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

func (ac *Client) AuthorizationURL(state string) string {
	return fmt.Sprintf("%s?redirect_uri=%s&client_id=%s&response_type=code&scope=%s&state=%s",
		ac.Provider.authURL,
		ac.RedirectURL,
		ac.Credentials.ClientID,
		strings.Join(ac.Scopes, " "),
		state,
	)
}

func (ac *Client) AccessToken(code string) (*Token, error) {
	u := fmt.Sprintf("%s?client_id=%s&client_secret=%s&code=%s",
		ac.Provider.tokenURL,
		ac.Credentials.ClientID,
		ac.Credentials.ClientSecret,
		code,
	)
	t := Token{}
	c := resty.New()
	resp, err := c.R().
		EnableTrace().
		SetResult(&t).
		SetHeader("Accept", "application/json").
		Post(u)
	if err != nil {
		return nil, fmt.Errorf("get access token: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode())
	}

	// TODO(rg): Handle generic error response here. Can happen if the secret is incorrect, for example.

	if t.AccessToken == "" {
		return nil, fmt.Errorf("empty access token")
	}
	return &t, nil
}

type Option func(client *Client)

func WithScopes(scopes []string) Option {
	return func(client *Client) {
		client.Scopes = scopes
	}
}
