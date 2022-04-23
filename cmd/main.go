package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/google/go-github/v43/github"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"gofuzzyclone/internal/logger"
	"gofuzzyclone/internal/prompter"
	githubService "gofuzzyclone/pkg/github"

	"github.com/briandowns/spinner"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

var configFilePath = os.ExpandEnv("$HOME/.gofuzzyclone.json")

type config struct {
	Ghtoken string `json:"github_token"`
}

func (c *cli) renewGhToken() (ghToken string) {
	c.prompter.Highlight("Generate token at https://github.com/settings/tokens/new?scopes=repo&description=gofuzzyclone-cli")
	c.prompter.Gather("GITHUB_TOKEN:", &ghToken)
	err := c.github.ValidateToken(ghToken)
	if err != nil {
		c.logger.HandleError(errors.Wrap(err, "invalid token"))
	} else {
		c.prompter.Success("Token is valid")
	}
	return ghToken
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

type cli struct {
	github     *githubService.GithubService
	config     config
	ctx        context.Context
	mode       string
	owner      string
	search     string
	outputPath string
	logger     *logger.Logger
	prompter   *prompter.Prompter
}

func (c *cli) getlocalConfig() (cfg config, err error) {
	if _, err := os.Stat(configFilePath); err == nil {
		// read config file
		jsonBlob, err := ioutil.ReadFile(configFilePath)
		c.logger.HandleError(err)
		err = json.Unmarshal(jsonBlob, &cfg)
		c.logger.HandleError(err)
		return cfg, nil
	} else {
		// create config file
		token := c.renewGhToken()
		f, err := os.Create(configFilePath)
		c.logger.HandleError(err)
		defer func(f *os.File) {
			err := f.Close()
			c.logger.HandleError(err)
		}(f)
		cfg = config{
			Ghtoken: token,
		}
		jsonBlob, err := json.Marshal(cfg)
		c.logger.HandleError(err)
		err = ioutil.WriteFile(configFilePath, jsonBlob, 0644)
		c.logger.HandleError(err)
		return config{}, nil
	}
}

func main() {
	var setup = flag.Bool("setup", false, "renew github token")
	var outputPath = flag.String("output", "", "output to which directory")
	var search = flag.String("search", "", "search patern")
	var owner = flag.String("owner", "", "github user/org")
	var mode = flag.String("mode", "regex", "matching mechanism")
	flag.Parse()

	cli := cli{
		mode:       *mode,
		outputPath: *outputPath,
		owner:      *owner,
		search:     *search,
		ctx:        context.Background(),
		logger:     logger.New(),
		prompter:   prompter.New(),
	}
	if *setup {
		cli.renewGhToken()
		os.Exit(0)
	}
	if os.Getenv("CI") == "true" {
		if ghToken := os.Getenv("GITHUB_TOKEN"); ghToken != "" {
			cli.config.Ghtoken = ghToken
		} else {
			cli.logger.HandleError(errors.New("GITHUB_TOKEN is not set"))
		}
	} else {
		// get local config
		cfg, err := cli.getlocalConfig()
		cli.logger.HandleError(err)
		cli.config = cfg
	}

	svc := githubService.New(cli.ctx, cli.config.Ghtoken)
	cli.github = svc
	// check critical arguments from user
	if cli.search == "" {
		cli.prompter.Gather("Search:", &cli.search)
	}
	if cli.owner == "" {
		cli.prompter.Gather("Owner:", &cli.owner)
	}

	// let's go search for repos
	s := spinner.New(spinner.CharSets[4], 200*time.Millisecond) // Build our new spinner
	s.Prefix = fmt.Sprintf("Searching %q under %q ", cli.search, cli.owner)
	s.Start() // Start the spinner
	user, _ := cli.github.GetUser(cli.owner)

	var err error
	var allRepos []*github.Repository
	if *user.Type == "Organization" {
		allRepos, err = cli.github.GetOrgRepos(cli.owner)
	} else {
		allRepos, err = cli.github.GetPersonalRepos(cli.owner)
	}
	cli.logger.HandleError(err)

	allReposMatched := make([]*github.Repository, 0)
	for _, repo := range allRepos {
		matched := false
		if cli.mode == "regex" {
			matched, _ = regexp.MatchString(cli.search, repo.GetName())
		} else if cli.mode == "wildcard" {
			matched, _ = regexp.MatchString(wildCardToRegexp(cli.search), repo.GetName())
		}
		if matched {
			allReposMatched = append(allReposMatched, repo)
		}
	}
	fmt.Println("")
	s.Stop() // Stop the spinner when we're done
	for i, repo := range allReposMatched {
		fmt.Printf("%d %v\n", i+1, repo.GetFullName())
	}
	confirm := "n"

	cli.logger.Success(fmt.Sprintf("Result: %d repos match %q in %v", len(allReposMatched), cli.search, cli.owner))
	if len(allReposMatched) == 0 {
		os.Exit(0)
	}

	if cli.outputPath == "" {
		cli.prompter.Gather(fmt.Sprintf("Clone %d repos to which folder?", len(allReposMatched)), &cli.outputPath)
	}

	cli.prompter.Highlight("Confirm? (Y/n)")
	cli.prompter.Gather("", &confirm)

	if confirm != "Y" {
		cli.logger.Error("Aborted")
		os.Exit(0)
	}

	for _, repo := range allReposMatched {
		_, err := git.PlainClone(path.Join(cli.outputPath, repo.GetName()), false, &git.CloneOptions{
			URL: repo.GetCloneURL(),
			Auth: &http.BasicAuth{
				Username: user.GetName(),
				Password: cli.config.Ghtoken,
			},
			Depth: 1,
		})
		if err != nil {
			if err.Error() == "repository already exists" {
				cli.logger.Warn(fmt.Sprintf("Skipped: %v (already exists)", repo.GetName()))
			} else {
				cli.logger.Error(fmt.Sprintf("Failed: %v", err.Error()))
			}
		} else {
			cli.logger.Success(fmt.Sprintf("Cloned: %v", repo.GetName()))
		}
	}

	cli.logger.Info(fmt.Sprintf("DONE with %d repositories", len(allReposMatched)))
}
