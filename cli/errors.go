package cli

import (
	"errors"

	"github.com/tamnd/gitee-cli/gitee"
)

func isNotFound(err error) bool {
	return errors.Is(err, gitee.ErrNotFound)
}
