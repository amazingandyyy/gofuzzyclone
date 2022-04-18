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

	"gofuzzyclone/internal/helper"
	"gofuzzyclone/internal/logger"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/google/go-github/v43/github"
	"golang.org/x/oauth2"
)

var configFilePath string = os.ExpandEnv("$HOME/.gofuzzyclone.json")

type config struct {
	Ghtoken string `json:"github_token"`
}

func validateGhToken(ghToken string) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: ghToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	_, _, err := client.Repositories.List(ctx, "github", nil)
	helper.HandleError(err, "Invalid token")
}

func renewGhToken() {
	ghToken := ""
	hasConfigFile := false
	logger.Println("yellow", "Generate token at https://github.com/settings/tokens/new?scopes=repo&description=gofuzzyclone-cli")
	fmt.Println("GITHUB_TOKEN: ")
	fmt.Scanf("%s", &ghToken)
	validateGhToken(ghToken)
	logger.Println("green", "Token is valid")

	var jsonBlob = []byte(`{"github_token": "` + ghToken + `"}`)
	cf := config{}
	err := json.Unmarshal(jsonBlob, &cf)
	helper.HandleError(err)
	if _, err := os.Stat(configFilePath); err == nil {
		hasConfigFile = true
	}
	if os.Getenv("CI") != "true" {
		cfJson, _ := json.Marshal(cf)
		if !hasConfigFile {
			// create one
			f, err := os.Create(configFilePath)
			helper.HandleError(err)
			defer f.Close()
			err = ioutil.WriteFile(configFilePath, cfJson, 0644)
			helper.HandleError(err)
		}
	}
}

// https://github.com/vrenjith/github-pr-manager/blob/5120424b4a9f4ac4675080eebce8582c5905e626/github.go#L20
func getOrgRepos(client *github.Client, org string) ([]*github.Repository, error) {
	max := 1000
	var allRepos []*github.Repository
	opt := &github.RepositoryListByOrgOptions{
		Sort: "pushed",
		ListOptions: github.ListOptions{
			PerPage: 80,
		},
	}
	counter := 0
	dotsCounter := 0
	dots := []string{".", ".", ".", "."}
	for {
		fmt.Printf("\rScanning %d repositories %s", counter, strings.Join(dots[:dotsCounter%3+1], ""))
		dotsCounter++
		repos, resp, err := client.Repositories.ListByOrg(context.Background(), org, opt)
		counter += len(repos)
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
		Sort:       "pushed",
		Visibility: "all",
		ListOptions: github.ListOptions{
			PerPage: 60,
		},
	}
	counter := 0
	dotsCounter := 0
	dots := []string{".", ".", ".", ".", "."}
	for {
		fmt.Printf("\rScanning %d repositories %s", counter, strings.Join(dots[:dotsCounter%4], " "))
		dotsCounter++
		repos, resp, err := client.Repositories.List(context.Background(), owner, opt)
		counter += len(repos)
		if len(allRepos) > max {
			allRepos = allRepos[:max]
			break
		}
		helper.HandleError(err, fmt.Sprintf("Unable to fetch respositories for the user: %s", err))
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
	if userType == "Organization" {
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
	var dest = flag.String("output", "", "output to which directory")
	var search = flag.String("search", "", "search patern")
	var owner = flag.String("owner", "", "github user/org")
	var mode = flag.String("mode", "regex", "matching mechanism")
	flag.Parse()

	if *setup {
		renewGhToken()
		os.Exit(0)
	}
	ghToken := os.Getenv("GITHUB_TOKEN")
	if _, err := os.Stat(configFilePath); err == nil {
		// config file does exist
		existingConfig, _ := os.ReadFile(configFilePath)
		var cf config
		json.Unmarshal(existingConfig, &cf)
		if cf.Ghtoken != "" {
			ghToken = cf.Ghtoken
		}
		validateGhToken(ghToken)
	}

	if ghToken == "" {
		renewGhToken()
	}

	if *search == "" {
		fmt.Printf("Search: ")
		fmt.Scanf("%s", search)
	}
	if *owner == "" {
		fmt.Printf("Owner: ")
		fmt.Scanf("%s", owner)
	}

	// s := spinner.New(spinner.CharSets[4], 200*time.Millisecond) // Build our new spinner
	// s.Prefix = fmt.Sprintf("Searching %q ", *search)
	// s.Start() // Start the spinner

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: ghToken})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	user, _, _ := client.Users.Get(ctx, *owner)

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
	// s.Stop()
	for i, repo := range all_repos_matched {
		fmt.Printf("%d %v\n", i+1, repo.GetFullName())
	}
	confirm := "n"
	logger.Println("green", fmt.Sprintf("Result: %d repos match %q in %v", len(all_repos_matched), *search, *owner))
	if len(all_repos_matched) == 0 {
		os.Exit(0)
	}

	if *dest == "" {
		fmt.Printf("Clone %d repos to which folder? ", len(all_repos_matched))
		fmt.Scanf("%s", dest)
	}

	logger.Printf("yellow", "Confirm? (Y/n) ")
	fmt.Scanf("%s", &confirm)

	if confirm != "Y" {
		logger.Println("red", "Aborted")
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
				logger.Println("yellow", fmt.Sprintf("Skipped: %v (already exists)", repo.GetName()))
			} else {
				logger.Println("red", fmt.Sprintf("Error cloning: %v", err.Error()))
			}
		} else {
			logger.Println("green", fmt.Sprintf("Cloned: %v", repo.GetName()))
		}
	}

	logger.Println("green", fmt.Sprintf("DONE with %d repositories", len(all_repos_matched)))
}
