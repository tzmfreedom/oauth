package provider

func init() {
	registerEndpoint("yahoo", &Endpoint{
		AuthURL:  "https://auth.login.yahoo.co.jp/yconnect/v1/authorization",
		TokenURL: "https://auth.login.yahoo.co.jp/yconnect/v1/token",
	})
}
