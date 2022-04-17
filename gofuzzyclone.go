package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v43/github"
	"golang.org/x/oauth2"
)

var gofuzzycloneConfigFilePath string = os.ExpandEnv("$HOME/.gofuzzyclone.json")

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func println(color string, message string) {
	color = strings.ToLower(color)
	var color_palettes = map[string]string{
		"red":    "\033[31m",
		"green":  "\033[32m",
		"yellow": "\033[33m",
		"reset":  "\033[0m",
	}

	fmt.Printf("%s%s%s\n", string(color_palettes[color]), message, color_palettes["reset"])
}

func printf(color string, message string) {
	color = strings.ToLower(color)
	var color_palettes = map[string]string{
		"red":    "\033[31m",
		"green":  "\033[32m",
		"yellow": "\033[33m",
		"reset":  "\033[0m",
	}

	fmt.Printf("%s%s%s", string(color_palettes[color]), message, color_palettes["reset"])
}

type gofuzzycloneConfig struct {
	Ghtoken string `json:"github_token"`
}

func validateGhToken(ghToken string) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: ghToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	_, _, err := client.Repositories.List(ctx, "github", nil)
	if err != nil {
		println("red", "Invalid token")
		os.Exit(1)
	}
}

func renewGhToken() {
	ghToken := ""
	hasConfigFile := false
	println("yellow", "Generate token at https://github.com/settings/tokens/new?scopes=repo&description=gofuzzyclone-cli")
	fmt.Println("GITHUB_TOKEN: ")
	fmt.Scanf("%s", &ghToken)
	validateGhToken(ghToken)
	println("green", "Token is valid")

	var jsonBlob = []byte(`{"github_token": "` + ghToken + `"}`)
	gofuzzycloneConfig := gofuzzycloneConfig{}
	err := json.Unmarshal(jsonBlob, &gofuzzycloneConfig)
	if err != nil {
		panic(err)
	}
	if _, err := os.Stat(gofuzzycloneConfigFilePath); err == nil {
		hasConfigFile = true
	}
	if os.Getenv("CI") != "true" {
		gofuzzycloneConfigJson, _ := json.Marshal(gofuzzycloneConfig)
		if !hasConfigFile {
			// create one
			f, err := os.Create(gofuzzycloneConfigFilePath)
			check(err)
			defer f.Close()
			err = ioutil.WriteFile(gofuzzycloneConfigFilePath, gofuzzycloneConfigJson, 0644)
			check(err)
		}
	}
}

// https://github.com/vrenjith/github-pr-manager/blob/5120424b4a9f4ac4675080eebce8582c5905e626/github.go#L20
func getOrgRepos(client *github.Client, org string) ([]*github.Repository, error) {
	max := 1000
	var allRepos []*github.Repository
	opt := &github.RepositoryListByOrgOptions{
		Sort: "updated",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	for {
		repos, resp, err := client.Repositories.ListByOrg(context.Background(), org, opt)
		if len(allRepos) > max {
			allRepos = allRepos[:max]
			break
		}
		if err != nil {
			log.Fatal("Unable to fetch respositories for the organization: ", err)
			return nil, err
		}
		if len(repos) == 0 {
			break
		}

		opt.Page = resp.NextPage
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
	}
	return allRepos, nil
}

func getPersonalRepos(client *github.Client, owner string) ([]*github.Repository, error) {
	max := 1000
	var allRepos []*github.Repository
	opt := &github.RepositoryListOptions{
		Sort: "updated",
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}
	for {
		repos, resp, err := client.Repositories.List(context.Background(), owner, opt)
		if len(allRepos) > max {
			allRepos = allRepos[:max]
			break
		}
		if err != nil {
			log.Fatal("Unable to fetch respositories for the user: ", err)
			return nil, err
		}
		if len(repos) == 0 {
			break
		}

		opt.Page = resp.NextPage
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
	}
	return allRepos, nil
}

func getRepos(client *github.Client, userType, owner string) ([]*github.Repository, error) {
	if userType == "org" {
		return getOrgRepos(client, owner)
	} else {
		return getPersonalRepos(client, owner)
	}
}

// wildCardToRegexp converts a wildcard pattern to a regular expression pattern.
func wildCardToRegexp(pattern string) string {
	var result strings.Builder
	for i, literal := range strings.Split(pattern, "*") {

		// Replace * with .*
		if i > 0 {
			result.WriteString(".*")
		}

		// Quote any regular expression meta characters in the
		// literal text.
		result.WriteString(regexp.QuoteMeta(literal))
	}
	return result.String()
}

func main() {
	var setup = flag.Bool("setup", false, "renew github token")
	var dest = flag.String("out", ".", "output to which sdirectory")
	var search = flag.String("search", "", "search patern")
	var owner = flag.String("owner", "", "github user/org")
	var mode = flag.String("mode", "regex", "matching mechanism")
	flag.Parse()

	if *setup {
		renewGhToken()
		os.Exit(0)
	}
	ghToken := os.Getenv("GITHUB_TOKEN")
	if _, err := os.Stat(gofuzzycloneConfigFilePath); err == nil {
		// config file does exist
		config, _ := os.ReadFile(gofuzzycloneConfigFilePath)
		var gofuzzycloneConfig gofuzzycloneConfig
		err := json.Unmarshal(config, &gofuzzycloneConfig)
		if err != nil {
			panic(err)
		}
		if gofuzzycloneConfig.Ghtoken != "" {
			ghToken = gofuzzycloneConfig.Ghtoken
		}
		validateGhToken(ghToken)
	}

	if ghToken == "" {
		renewGhToken()
	}

	if *search == "" {
		fmt.Printf("search: ")
		fmt.Scanf("%s", search)
	}
	if *owner == "" {
		fmt.Printf("owner: ")
		fmt.Scanf("%s", owner)
	}

	s := spinner.New(spinner.CharSets[39], 200*time.Millisecond) // Build our new spinner
	s.Prefix = fmt.Sprintf("Searching %q ", *search)
	s.Start() // Start the spinner

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: ghToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	user, _, _ := client.Users.Get(ctx, "")

	all_repos, _ := getRepos(client, *user.Type, *owner)
	all_repos_matched := make([]*github.Repository, 0)

	for _, repo := range all_repos {
		matched := false
		if *mode == "regex" {
			matched, _ = regexp.MatchString(*search, repo.GetName())
		} else if *mode == "wildcard" {
			matched, _ = regexp.MatchString(wildCardToRegexp(*search), repo.GetName())
		}
		if matched {
			all_repos_matched = append(all_repos_matched, repo)
		}
	}
	fmt.Println("")
	s.Stop()
	for i, repo := range all_repos_matched {
		fmt.Printf("%d %v\n", i+1, repo.GetFullName())
	}
	confirm := "n"
	println("green", fmt.Sprintf("Result: %d repos match %q", len(all_repos_matched), *search))

	if *dest == "./" {
		fmt.Printf("Clone all repos to which folder? ")
		fmt.Scanf("%s", dest)
	}

	printf("yellow", fmt.Sprintf("Clone %d repos to %v? (Y/n) ", len(all_repos_matched), *dest))
	fmt.Scanf("%s", &confirm)

	if confirm != "Y" {
		println("red", "Aborted")
		os.Exit(0)
	}

	for _, repo := range all_repos_matched {
		_, err := git.PlainClone(path.Join(*dest, repo.GetName()), false, &git.CloneOptions{
			URL: repo.GetCloneURL(),
			Auth: &http.BasicAuth{
				Username: user.GetName(),
				Password: ghToken,
			},
			Depth: 1,
		})
		if err != nil {
			if err.Error() == "repository already exists" {
				println("yellow", fmt.Sprintf("Skipped: %v (already exists)", repo.GetName()))
			} else {
				println("red", fmt.Sprintf("Error cloning: %v", err.Error()))
			}
		} else {
			println("green", fmt.Sprintf("Cloned: %v", repo.GetName()))
		}
	}

	println("green", fmt.Sprintf("DONE with %d repositories", len(all_repos_matched)))
}
