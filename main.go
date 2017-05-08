package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/skratchdot/open-golang/open"
	"github.com/tzmfreedom/oauth/provider"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
)

var (
	Version  string
	Revision string
)

func main() {
	cli.VersionPrinter = func(c *cli.Context) {
		fmt.Printf("version=%s revision=%s\n", c.App.Version, Revision)
	}

	app := cli.NewApp()
	app.Name = "oauth"
	app.Usage = "oauth command line tool"
	app.Version = Version
	app.Commands = []cli.Command{
		{
			Name:    "authorize",
			Aliases: []string{"a"},
			Usage:   "authorize command",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "authorize_url",
				},
				cli.StringFlag{
					Name: "token_url",
				},
				cli.StringFlag{
					Name:  "response_type",
					Value: "code",
				},
				cli.StringFlag{
					Name: "client_id",
				},
				cli.StringFlag{
					Name: "client_secret",
				},
				cli.StringFlag{
					Name: "redirect_uri",
				},
				cli.StringFlag{
					Name: "scope",
				},
				cli.StringFlag{
					Name: "state",
				},
				cli.BoolFlag{
					Name: "random_state",
				},
				cli.BoolFlag{
					Name: "interactive",
				},
				cli.BoolFlag{
					Name: "open",
				},
				cli.BoolFlag{
					Name: "auto",
				},
				cli.IntFlag{
					Name:  "port",
					Value: 1234,
				},
				cli.StringFlag{
					Name: "provider",
				},
			},
			Action: func(c *cli.Context) error {
				authConfig, err := getOAuthConfig(c)
				if err != nil {
					return err
				}
				if authConfig.Endpoint.AuthURL == "" {
					return errors.New("authorization_url rquired")
				}
				state := getState(c)
				switch getResponseType(c) {
				case "code":
					if c.Bool("auto") {
						if authConfig.Endpoint.TokenURL == "" {
							return errors.New("token_url rquired")
						}
						oauthDanceCodeGrant(state, authConfig)
						return nil
					}
					authURL := authConfig.AuthCodeURL(state)
					if c.Bool("open") {
						open.Start(authURL)
					} else {
						fmt.Println(authURL)
					}
				case "token":
					if c.Bool("auto") {
						oauthDanceImplicitGrant(state, authConfig)
						return nil
					}
					authURL := authConfig.AuthCodeURL(state, oauth2.SetAuthURLParam("response_type", "token"))
					if c.Bool("open") {
						open.Start(authURL)
					} else {
						fmt.Println(authURL)
					}
				}
				return nil
			},
		},
		{
			Name:    "token",
			Aliases: []string{"a"},
			Usage:   "get token command",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "token_url",
				},
				cli.StringFlag{
					Name:  "grant_type",
					Value: "authorization_code",
				},
				cli.StringFlag{
					Name: "client_id",
				},
				cli.StringFlag{
					Name: "client_secret",
				},
				cli.StringFlag{
					Name: "redirect_uri",
				},
				cli.StringFlag{
					Name: "scope",
				},
				cli.StringFlag{
					Name: "code",
				},
				cli.BoolFlag{
					Name: "interactive",
				},
			},
			Action: func(c *cli.Context) error {
				authConfig, err := getOAuthConfig(c)
				if err != nil {
					return err
				}
				if authConfig.Endpoint.TokenURL == "" {
					return errors.New("token_url required")
				}
				var token *oauth2.Token
				switch getGrantType(c) {
				case "code":
					code := getCode(c)
					if code == "" {
						return errors.New("code required")
					}
					token, err = authConfig.Exchange(nil, code)
					if err != nil {
						return err
					}
				case "password":
					username, password := getPasswordCredentials(c)
					if username == "" || password == "" {
						return errors.New("credentials required")
					}
					token, err = authConfig.PasswordCredentialsToken(nil, username, password)
					if err != nil {
						return err
					}
				}
				fmt.Println("AccessToken\t%s", token.AccessToken)
				fmt.Println("RefreshToken\t%s", token.RefreshToken)
				fmt.Println("TokenType\t%s", token.TokenType)
				fmt.Println("Expiry\t%s", token.Expiry)
				return nil
			},
		},
	}
	app.Run(os.Args)
}

func ask(question, defaultAnswer string) string {
	fmt.Printf("%s [%s]:", question, defaultAnswer)
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	answer = strings.TrimRight(answer, "\n")
	if answer == "" {
		return defaultAnswer
	}
	return answer
}
func getResponseType(c *cli.Context) *oauth2.Config {
	responseType := c.String("response_type")
	if c.Bool("interactive") {
		responseType = ask("response_type", responseType)
	}
	return responseType
}

func getOAuthConfig(c *cli.Context) (*oauth2.Config, errror) {
	var authorizeUrl, tokenUrl string
	switch c.String("provider") {
	case "salesforce":
		authorizeUrl = provider.Salesforce.AuthURL
		tokenUrl = provider.Salesforce.TokenURL
	default:
		authorizeUrl = c.String("authorize_url")
		tokenUrl = c.String("token_url")

	}
	clientId := c.String("client_id")
	clientSecret := c.String("client_secret")
	redirectUri := c.String("redirect_uri")
	if c.Bool("auto") {
		redirectUri = fmt.Sprintf("http://localhost:%d", c.Int("port"))
	}
	scope := c.String("scope")
	if c.Bool("interactive") {
		authorizeUrl = ask("authorize_url", authorizeUrl)
		tokenUrl = ask("token_url", tokenUrl)
		clientId = ask("client_id", clientId)
		clientSecret = ask("client_secret", clientSecret)
		redirectUri = ask("redirect_uri", redirectUri)
		scope = ask("scope", scope)
	}
	if clientId == "" {
		return nil, errors.New("client_id required")
	}
	if clientSecret == "" {
		return nil, errors.New("client_secret required")
	}
	if redirectUri == "" {
		return nil, errors.New("redirect_uri required")
	}
	return &oauth2.Config{
		Endpoint:     oauth2.Endpoint{TokenURL: tokenUrl, AuthURL: authorizeUrl},
		ClientID:     clientId,
		ClientSecret: clientSecret,
		RedirectURL:  redirectUri,
		Scopes:       strings.Split(scope, ","),
	}, nil
}

func getState(c *cli.Context) string {
	state := c.String("state")
	if c.Bool("interactive") {
		state = ask("state", state)
	}
	if c.Bool("random_state") {
		b := make([]byte, 32)
		rand.Read(b)
		state = base64.URLEncoding.EncodeToString(b)
	}
	return state
}

func oauthDanceCodeGrant(state string, c *oauth2.Config) {
	receive := make(chan string)
	s := &http.Server{
		Addr: fmt.Sprintf(":%d", c.Int("port")),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			receive <- r.URL.Query().Get("code")
			w.Write([]byte("<script>window.open('about:blank','_self').close()</script>"))
			w.(http.Flusher).Flush()
		}),
	}
	go s.ListenAndServe()

	open.Start(c.AuthCodeURL(state, nil))
	code := <-receive
	token, err := c.Exchange(context.Background(), code)
	if err != nil {
		return err
	}
	fmt.Println(token.AccessToken)
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	s.Shutdown(ctx)
}

func oauthDanceImplicitGrant(state string, c *oauth2.Config) {
	receive := make(chan string)
	s := &http.Server{
		Addr: fmt.Sprintf(":%d", c.Int("port")),
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				w.Write([]byte("<script>location.href = '/close?' + location.hash.substring(1);</script>"))
				w.(http.Flusher).Flush()
			} else {
				w.Write([]byte("<html><body>Close this window.</body></html>"))
				w.(http.Flusher).Flush()
				receive <- r.URL.Query().Get("access_token")
			}
		}),
	}
	go s.ListenAndServe()

	open.Start(c.AuthCodeURL(state, oauth2.SetAuthURLParam("response_type", "token")))
	accessToken := <-receive
	fmt.Println(accessToken)
	ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
	s.Shutdown(ctx)
}

func getGrantType(c *cli.Context) string {
	grantType := c.String("grant_type")
	if c.Bool("interactive") {
		grantType = ask("grant_type", "authorization_code")
	}
	return grantType
}

func getPasswordCredentials(c *cli.Context) (username, password) {
	username := c.String("username")
	password := c.String("password")
	if c.Bool("interactive") {
		username = ask("username", username)
		password = ask("password", password)
	}
	return username, password

}

func getCode(c *cli.Context) string {
	code := c.String("code")
	if c.Bool("interactive") {
		code = ask("code", "")
	}
	return code
}
