// Package cli assembles the gitee command tree from the gitee domain on top of
// the any-cli/kit framework.
package cli

import (
	"github.com/tamnd/any-cli/kit"
	"github.com/tamnd/gitee-cli/gitee"
)

// Build metadata, set via -ldflags at release time.
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// NewApp assembles the kit application from the gitee domain.
func NewApp() *kit.App {
	id := gitee.BaseIdentity()
	id.Version = Version
	app := kit.New(id, kit.WithDefaults(gitee.Defaults))
	gitee.Register(app)
	return app
}
