# OAuth Command

Simple CLI to request OAuth(WIP)

## Install

For Linux and macOS user
```bash
$ curl -sL http://install.freedom-man.com/oauth | bash
```

If you want to install zsh completion, add --zsh-completion option
```bash
$ curl -sL http://install.freedom-man.com/oauth | bash -s -- --zsh-completion
```

For Windows user
```bash
$ (New-Object Net.WebClient).DownloadString('http://install.freedom-man.com/oauth.ps1') | iex
```

If you want to use latest version, execute following command.
```bash
$ go get -u github.com/tzmfreedom/oauth
```

## Usage

```
NAME:
   oauth - oauth command line tool

USAGE:
   oauth [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     authorize, a  authorize command
     token, a      get token command
     help, h       Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

### Authorize

```
ME:
   oauth authorize - authorize command

USAGE:
   oauth authorize [command options] [arguments...]

OPTIONS:
   --authorize_url value
   --token_url value
   --response_type value       (default: "code")
   --client_id value
   --client_secret value
   --redirect_uri value
   --scope value
   --state value
   --random_state
   --interactive, -I
   --open, -O
   --auto, -A
   --port value                (default: 1234)
   --provider value, -P value
```

Authorize subcommand displays oauth2.0 authorize url by default.
```bash
$ oauth authorize --client_id {client_id} --client_secret {client_secret} \
  --redirect_uri {uri} --scope {scope} --state {state} --authorize_url {uri}
```
The parameters of client_id, client_secret, redirect_uri is required.  
You can use state_random option instead of state option to get random token.

If you set interactive option, you can set oauth2.0 parameters interactively.
```bash
$ oauth authorize --interactive
```

If you want to get token automatically, you can set auto option.
```bash
$ oauth authorize --interactive --auto
```
If auto option is set, it display token response information.

You can set provider option instead of setting authorize_url and token_url.
There are example for google.
```bash
$ oauth authorize --interactive --provider google
```
Available value is google, facebook, yahoo, github, salesforce, slack and box.

### Token

```
NAME:
   oauth token - get token command

USAGE:
   oauth token [command options] [arguments...]

OPTIONS:
   --token_url value
   --grant_type value          (default: "authorization_code")
   --client_id value
   --client_secret value
   --redirect_uri value
   --scope value
   --code value
   --interactive, -I
   --provider value, -P value
```

Token subcommand displays oauth2.0 token response information.
```bash
$ oauth token --client_id {client_id} --client_secret {client_secret} \
  --redirect_uri {uri} --code {code} --state {state} --token_url {uri}
```

You can set interactive option in the same way as Authorize subcommand.
```bash
$ oauth token --interactive
```

## License

MIT
