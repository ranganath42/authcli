package authcli

type Provider struct {
	authURL  string
	tokenURL string
}

var (
	ProviderGitHub = Provider{
		authURL:  "https://github.com/login/oauth/authorize",
		tokenURL: "https://github.com/login/oauth/access_token",
	}
)
