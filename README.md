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

Authorize subcommand displays oauth2.0 authorize url by default.
```bash
$ oauth authorize --client_id {client_id} --client_secret {client_secret} --redirect_uri {uri} --scope {scope} --state {state}
```
The parameters of client_id, client_secret, redirect_uri is required.
You can use state_random option instead of state option to get random token

If you set interactive option, you can set oauth2.0 parameters interactively.
```bash
$ oauth authorize --interactive
```

If you want to get token automatically, you can set auto option.
```bash
$ oauth authorize --interactive --auto
```
If auto option is set, it display token response information.

Token subcommand displays oauth2.0 token response information.
```bash
$ oauth token --client_id {client_id} --client_secret {client_secret} --redirect_uri {uri} --code {code} --state {state}
```

You can set interactive option in the same way as Authorize subcommand.
```bash
$ oauth --interactive
```

## License

MIT
