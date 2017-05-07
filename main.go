package main

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/base64"
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
				var authorize_url, token_url string
				switch c.String("provider") {
				case "salesforce":
					authorize_url = provider.Salesforce.AuthURL
					token_url = provider.Salesforce.TokenURL
				default:
					authorize_url = c.String("authorize_url")
					token_url = c.String("token_url")

				}
				response_type := c.String("response_type")
				client_id := c.String("client_id")
				client_secret := c.String("client_secret")
				redirect_uri := c.String("redirect_uri")
				if c.Bool("auto") {
					redirect_uri = fmt.Sprintf("http://localhost:%d", c.Int("port"))
				}
				scope := c.String("scope")
				state := c.String("state")
				if c.Bool("interactive") {
					authorize_url = ask("authorize_url", authorize_url)
					token_url = ask("token_url", token_url)
					fmt.Println(response_type)
					response_type = ask("response_type", response_type)
					fmt.Println(response_type)
					client_id = ask("client_id", client_id)
					client_secret = ask("client_secret", client_secret)
					redirect_uri = ask("redirect_uri", redirect_uri)
					scope = ask("scope", scope)
					state = ask("state", state)
				}
				if c.Bool("random_state") {
					b := make([]byte, 32)
					rand.Read(b)
					state = base64.URLEncoding.EncodeToString(b)
				}
				authConfig := &oauth2.Config{
					Endpoint:     oauth2.Endpoint{AuthURL: authorize_url},
					ClientID:     client_id,
					ClientSecret: client_secret,
					RedirectURL:  redirect_uri,
					Scopes:       strings.Split(scope, ","),
				}
				switch response_type {
				case "code":
					authURL := authConfig.AuthCodeURL(state)
					if c.Bool("auto") {
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

						open.Start(authURL)
						code := <-receive
						authConfig := &oauth2.Config{
							Endpoint:     oauth2.Endpoint{TokenURL: token_url},
							ClientID:     client_id,
							ClientSecret: client_secret,
							RedirectURL:  redirect_uri,
						}
						token, err := authConfig.Exchange(context.Background(), code)
						if err != nil {
							return err
						}
						fmt.Println(token.AccessToken)
						ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
						s.Shutdown(ctx)
					} else if c.Bool("open") {
						open.Start(authURL)
					} else {
						fmt.Println(authURL)
					}
				case "token":
					authURL := authConfig.AuthCodeURL(state, oauth2.SetAuthURLParam("response_type", "token"))
					if c.Bool("auto") {
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

						open.Start(authURL)
						accessToken := <-receive
						fmt.Println(accessToken)
						ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
						s.Shutdown(ctx)
					} else if c.Bool("open") {
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
					Name: "endpoint",
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
					Name: "code",
				},
				cli.BoolFlag{
					Name: "interactive",
				},
			},
			Action: func(c *cli.Context) error {
				grant_type := c.String("grant_type")
				endpoint := c.String("endpoint")
				client_id := c.String("client_id")
				client_secret := c.String("client_secret")
				redirect_uri := c.String("redirect_uri")
				username := c.String("username")
				password := c.String("password")
				code := c.String("code")
				if c.Bool("interactive") {
					grant_type = ask("grant_type", "authorization_code")
					endpoint = ask("endpoint", "")
					client_id = ask("client_id", "")
					client_secret = ask("client_secret", "")
					redirect_uri = ask("redirect_uri", "")
					code = ask("code", "")
				}
				authConfig := &oauth2.Config{
					Endpoint:     oauth2.Endpoint{TokenURL: endpoint},
					ClientID:     client_id,
					ClientSecret: client_secret,
					RedirectURL:  redirect_uri,
				}
				var token *oauth2.Token
				var err error
				switch grant_type {
				case "code":
					token, err = authConfig.Exchange(nil, code)
					if err != nil {
						return err
					}
				case "password":
					token, err = authConfig.PasswordCredentialsToken(nil, username, password)
					if err != nil {
						return err
					}
				}
				fmt.Println(token.AccessToken)
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
