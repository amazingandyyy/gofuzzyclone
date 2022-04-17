# gofuzzyclone

go fuzzy search repos in regex and clone

# Installation

```sh
bash <(curl -sL https://raw.githubusercontent.com/amazingandyyy/gofuzzyclone/main/install.sh)
```

## Preparation

- [Generate a Github personal access token](https://github.com/settings/tokens/new?scopes=repo&description=gofuzzyclone-cli)
  - [repo] scrope
  - [no expiration]

## Usage

```
# get instructions
gofuzzyclone -help

# fastline mode
gofuzzyclone -owner amazingandyyy -search "^go.*" -out ./projects
gofuzzyclone -owner amazingandyyy -search "*-template" -mode wildcard -out ./projects

# interactive mode
gofuzzyclone
```

## Resources

- [regex](http://regex101.com)
- wildcard

## LICENSE

[MIT](./LICENSE)
