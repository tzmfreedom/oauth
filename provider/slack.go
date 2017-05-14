package provider

func init() {
	registerEndpoint("slack", &Endpoint{
		AuthURL:  "https://slack.com/oauth/authorize",
		TokenURL: "https://slack.com/api/oauth.access",
	})
}
