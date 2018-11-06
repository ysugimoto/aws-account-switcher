# aws-account-switcher

Command tool to set aws account to environment variable easily from credentials.

# Installation

Install via go runtime

```
$ go get github.com/ysugimoto/aws-acccount-switcher/...
```

Or, download prebuilt binary from Release page.

After that, you can use `acs` command.

# Usage

`acs` accepts profile name of first argument, if supply it, loads and parse `~/.aws/credentials` file and export as environment variable, so you can export it by using `source` command on osx/linux:

```
$ acs [profile name] | source -
```

If you don't supply any arguments, `acs` displays select interface:

```
$ acs | source -
default         : XXXXXXXXXXXX | YYYYYYYYYYYYYYYY
other-account   : AAAAAAAAAAAA | BBBBBBBBBBBBBBBB
...
```

You can choose account by cursor key or j(up)/k(down) key.

# License

MIT

# Author

Yoshiaki Sugimoto <sugimoto@wnotes.net>

