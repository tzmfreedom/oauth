package provider

func init() {
	registerEndpoint("facebook", &Endpoint{
		AuthURL:  "https://www.facebook.com/v2.9/dialog/oauth",
		TokenURL: "https://graph.facebook.com/v2.9/oauth/access_token",
	})
}
