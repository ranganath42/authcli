package main

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/ranganath42/authcli"
	"log"
	"net/http"
)

var token *authcli.Token

// rootHandler displays a link to the authorization URL of the provider.
// The user is sent to the provider's login page and asked to authenticate.
func rootHandler(url string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		htmlIndex := fmt.Sprintf("<html><body><a href='%s'>Login with GitHub</a></body></html>", url)
		fmt.Fprintf(w, htmlIndex)
	}
}

func callbackHandler(client *authcli.Client, state string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Printf("Failed to parse query: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if authErr := r.FormValue("error"); authErr != "" {
			errDesc := r.FormValue("error_description")
			log.Printf("Auth error: %s, description: %s", authErr, errDesc)
			w.WriteHeader(http.StatusBadRequest)
			htmlErr := fmt.Sprintf("<html><body><a>%s Back to </a><a href='/'>Home</a></body></html>", errDesc)
			fmt.Fprintf(w, htmlErr)
			return
		}

		if sts := r.FormValue("state"); sts != state {
			log.Printf("Invalid state")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		code := r.FormValue("code")
		if code == "" {
			log.Printf("Code not found")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		token, err = client.AccessToken(code)
		if err != nil {
			log.Printf("Failed to get access token: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		http.Redirect(w, r, "/welcome", http.StatusPermanentRedirect)
	}
}

func welcomeHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result := struct {
			Name      string `json:"name"`
			Location  string `json:"location"`
			AvatarURL string `json:"avatar_url"`
		}{}
		c := resty.New()
		resp, err := c.R().
			EnableTrace().
			SetHeader("Authorization", fmt.Sprintf("token %s", token.AccessToken)).
			SetResult(&result).
			Get("https://api.github.com/user")
		if err != nil {
			log.Printf("Failed to get user: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if resp.StatusCode() != http.StatusOK {
			log.Printf("Failed to get user: %v", resp)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		htmlWelcome := fmt.Sprintf("<html> <head> </head> <img src=%s width=\"64\" height=\"64\"><br>"+
			"<a>Welcome, %s. Back to </a><a href='/'>Home</a> </body> </html>",
			result.AvatarURL, result.Name,
		)
		fmt.Fprintf(w, htmlWelcome)
	}
}
