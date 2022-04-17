# gofuzzyclone

go fuzzy search repos and clone (supported wildcard/regex)

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
gofuzzyclone -help

gofuzzyclone
Search under which owner? amazingandyyy
Search for what wildcard pattern? (Press [Enter] to skip) go*
Search for what regex pattern? (Press [Enter] to skip) .*-template
Searching in amazingandyyy 🌍
1 amazingandyyy/goapprove
2 amazingandyyy/argocd-example-extension
3 amazingandyyy/learn-golang-with-stephen
4 amazingandyyy/algor-in-js
5 amazingandyyy/fotingo
6 amazingandyyy/good-job
7 amazingandyyy/go-grpc-start
8 amazingandyyy/gomorph
9 amazingandyyy/django-api
10 amazingandyyy/one-click-hugo-cms
11 amazingandyyy/javascript-algorithms
12 amazingandyyy/go-app
13 amazingandyyy/app-template
14 amazingandyyy/learn-golang-with-todd
15 amazingandyyy/learn-golang-basic
16 amazingandyyy/gotraining
17 amazingandyyy/go
18 amazingandyyy/goexpress
19 amazingandyyy/go-lang-cheat-sheet
20 amazingandyyy/mongo-stars
21 amazingandyyy/mongo-property
22 amazingandyyy/mongo-flashcard
Found 22 amazingandyyy's repositories match pattern of go* .*-template
Clone all repos to which folder? ./test             
Are you sure to continue cloning them all into ./test ? (Y/n) Y
Cloned: goapprove
Cloned: argocd-example-extension
Cloned: learn-golang-with-stephen
Cloned: algor-in-js
Cloned: fotingo
Cloned: good-job
Cloned: go-grpc-start
Cloned: gomorph
Cloned: django-api
Cloned: one-click-hugo-cms
Cloned: javascript-algorithms
Cloned: go-app
Cloned: app-template
Cloned: learn-golang-with-todd
Cloned: learn-golang-basic
Cloned: gotraining
...
```

## LICENSE

[MIT](LICENSE)
