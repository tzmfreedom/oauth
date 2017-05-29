package provider

func init() {
	registerEndpoint("box", &Endpoint{
		AuthURL:  "https://app.box.com/api/oauth2/authorize",
		TokenURL: "https://app.box.com/api/oauth2/token",
	})
}
