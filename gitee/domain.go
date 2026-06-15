package gitee

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/tamnd/any-cli/kit"
	"github.com/tamnd/any-cli/kit/errs"
)

func init() { kit.Register(Domain{}) }

type Domain struct{}

func (Domain) Info() kit.DomainInfo {
	return kit.DomainInfo{
		Scheme:   "gitee",
		Hosts:    []string{"gitee.com"},
		Identity: BaseIdentity(),
	}
}

func BaseIdentity() kit.Identity {
	return kit.Identity{
		Binary: "gitee",
		Short:  "A command-line for Gitee (码云) — China's Git hosting platform.",
		Long: `gitee reads public data from Gitee (码云) via the official REST API v5.

Browse users, repositories, commits, issues, pull requests, releases, and more.
All commands work without a token. Set GITEE_TOKEN for higher rate limits.

gitee is an independent tool and is not affiliated with Gitee or OSChina.`,
		Site: "https://gitee.com",
		Repo: "https://github.com/tamnd/gitee-cli",
	}
}

func Defaults(c *kit.Config) {
	d := DefaultConfig()
	c.Rate = d.Rate
	c.Timeout = d.Timeout
	c.Retries = d.Retries
	c.UserAgent = d.UserAgent
}

func (Domain) Register(app *kit.App) {
	app.SetClient(newClient)
	registerOps(app)
}

func Register(app *kit.App) { Domain{}.Register(app) }

type Session struct {
	Client *Client
	Quiet  bool
}

func (s *Session) Progressf(format string, args ...any) {
	if s == nil || s.Quiet {
		return
	}
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func newClient(_ context.Context, c kit.Config) (any, error) {
	cfg := DefaultConfig()
	if c.UserAgent != "" {
		cfg.UserAgent = c.UserAgent
	}
	if c.Rate > 0 {
		cfg.Rate = c.Rate
	}
	if c.Timeout > 0 {
		cfg.Timeout = c.Timeout
	}
	if c.Retries > 0 {
		cfg.Retries = c.Retries
	}
	cfg.Token = os.Getenv("GITEE_TOKEN")
	return &Session{Client: NewClient(cfg), Quiet: c.Quiet}, nil
}

func MapErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, ErrNotFound) {
		return errs.NotFound("%s", err)
	}
	if errors.Is(err, ErrUnauthorized) {
		return errs.Usage("authentication required — set GITEE_TOKEN")
	}
	return err
}

func (Domain) Classify(input string) (uriType, id string, err error) {
	return "", "", errs.Usage("unrecognised gitee reference: %q", input)
}

func (Domain) Locate(uriType, id string) (string, error) {
	return "", errs.Usage("gitee has no resource type %q", uriType)
}
