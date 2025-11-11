package policy

import (
	"github.com/mithileshgupta12/velaris/internal/middleware"
	"xorm.io/xorm"
)

type Policy interface {
	CanDelete(ctxUser middleware.CtxUser, id int64) bool
}

type Policies struct {
	BoardPolicy Policy
}

func InitPolicies(engine *xorm.Engine) *Policies {
	return &Policies{
		BoardPolicy: NewBoardPolicy(engine),
	}
}
