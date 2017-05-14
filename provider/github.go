package provider

func init() {
	registerEndpoint("github", &Endpoint{
		AuthURL:  "https://github.com/login/oauth/authorize",
		TokenURL: "https://github.com/login/oauth/access_token",
	})
}
