# gofuzzyclone

Go fuzzy search repos with regex or wildcard

# Installation

```sh
bash <(curl -sL https://raw.githubusercontent.com/amazingandyyy/gofuzzyclone/main/scripts/install.sh)
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
gofuzzyclone -owner amazingandyyy -search "^go.*" -output ./code
gofuzzyclone -owner amazingandyyy -search "*-template" -mode wildcard -output ./projects

# interactive mode
gofuzzyclone
```

## Resources

- [regex](http://regex101.com)
- wildcard

## LICENSE

[MIT](./LICENSE)
