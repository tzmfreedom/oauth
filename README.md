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

You can use docker to use spm command.
```bash
$ docker run --rm tzmfree/oauth install {REPO} u {USERNAME] -p {PASSWORD}

```

If you want to use latest version, execute following command.
```bash
$ go get -u github.com/tzmfreedom/oauth
```

For Windows user, use Linux virtual machine, such as docker or vagrant.
```bash
$ (New-Object Net.WebClient).DownloadString('http://install.freedom-man.com/oauth.ps1') | iex
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
