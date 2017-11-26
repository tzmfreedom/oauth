package provider

func init() {
  registerEndpoint("line", &Endpoint{
    AuthURL:  "https://access.line.me/oauth2/v2.1/authorize",
    TokenURL: "https://api.line.me/oauth2/v2.1/token",
  })
}
