package policy

import (
	"github.com/mithileshgupta12/velaris/internal/middleware"
	"xorm.io/xorm"
)

type Policy interface {
	CanView(ctxUser middleware.CtxUser, id int64) (bool, error)
	CanCreate(ctxUser middleware.CtxUser, id int64) (bool, error)
	CanUpdate(ctxUser middleware.CtxUser, id int64) (bool, error)
	CanDelete(ctxUser middleware.CtxUser, id int64) (bool, error)
}

type Policies struct {
	BoardPolicy Policy
	ListPolicy  Policy
}

func InitPolicies(engine *xorm.Engine) *Policies {
	return &Policies{
		BoardPolicy: NewBoardPolicy(engine),
		ListPolicy:  NewListPolicy(engine),
	}
}
