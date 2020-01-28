package grifts

import (
	"github.com/slatunje/k8sroles/actions"

	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
