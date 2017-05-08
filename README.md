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


```bash
$ oauth authorize --client_id {client_id} --client_secret {client_secret} --redirect_uri {uri} --scope {scope} --state {state}
$ oauth authorize --client_id {client_id} --client_secret {client_secret} --redirect_uri {uri} --scope {scope} --state_random
```

--open

token
```bash
$ oauth token --client_id {client_id} --client_secret {client_secret} --redirect_uri {uri} --code {code} --state {state}
```

```bash
$ oauth --interactive
```

## License

MIT
