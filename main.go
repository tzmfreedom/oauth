package main

import (
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/urfave/cli"
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
					Name: "response_type",
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
			},
			Action: func(c *cli.Context) error {
				authorize_url := c.String("authorize_url")
				response_type := c.String("response_type")
				client_id := c.String("client_id")
				client_secret := c.String("client_secret")
				redirect_uri := c.String("redirect_uri")
				scope := c.String("scope")
				state := c.String("state")
				if c.Bool("interactive") {
					authorize_url = ask("authorize_url", "")
					response_type = ask("response_type", "code")
					client_id = ask("client_id", "")
					client_secret = ask("client_secret", "")
					redirect_uri = ask("redirect_uri", "")
					scope = ask("scope", "")
					state = ask("state", "")
				}
				if c.Bool("random_state") {
					b := make([]byte, 32)
					rand.Read(b)
					state = base64.URLEncoding.EncodeToString(b)
				}
				auth_url := createAuthorizeUrl(
					authorize_url,
					response_type,
					client_id,
					client_secret,
					redirect_uri,
					scope,
					state,
				)
				if c.Bool("open") {
					if err := exec.Command("open", auth_url).Run(); err != nil {
						return err
					}
				} else {
					fmt.Println(auth_url)
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
					Name: "grant_type",
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
				token_url := c.String("token_url")
				client_id := c.String("client_id")
				client_secret := c.String("client_secret")
				redirect_uri := c.String("redirect_uri")
				code := c.String("code")
				if c.Bool("interactive") {
					grant_type = ask("grant_type", "authorization_code")
					token_url = ask("token_url", "")
					client_id = ask("client_id", "")
					client_secret = ask("client_secret", "")
					redirect_uri = ask("redirect_uri", "")
					code = ask("code", "")
				}
				r, err := requestToken(
					grant_type,
					token_url,
					client_id,
					client_secret,
					redirect_uri,
					code,
				)
				if err != nil {
					return err
				}
				body, err := ioutil.ReadAll(r.Body)
				if err != nil {
					return err
				}
				fmt.Println(body)
				return nil
			},
		},
	}
	app.Run(os.Args)
}

func createAuthorizeUrl(authorize_url, response_type, client_id, client_secret, redirect_uri, scope, state string) string {
	v := url.Values{}
	v.Add("response_type", response_type)
	v.Add("client_id", client_id)
	v.Add("client_secret", client_secret)
	v.Add("redirect_uri", redirect_uri)
	if scope != "" {
		v.Add("scope", scope)
	}
	if state != "" {
		v.Add("state", state)
	}

	return authorize_url + "?" + v.Encode()
}

func ask(question, defaultAnswer string) string {
	fmt.Print(question + " : ")
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	if answer == "" {
		return defaultAnswer
	}
	return answer
}

func requestToken(grant_type, token_url, client_id, client_secret, redirect_uri string, code) (*http.Response, error) {
	v := url.Values{}
	v.Add("grant_type", grant_type)
	v.Add("client_id", client_id)
	v.Add("client_secret", client_secret)
	v.Add("redirect_uri", redirect_uri)
	v.Add("code", code)

	req, _ := http.NewRequest("POST", token_url, strings.NewReader(v.Encode()))
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
