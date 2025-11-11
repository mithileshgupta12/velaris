package policy

import (
	"github.com/mithileshgupta12/velaris/internal/middleware"
	"xorm.io/xorm"
)

type boardPolicy struct {
	engine *xorm.Engine
}

func NewBoardPolicy(engine *xorm.Engine) Policy {
	return &boardPolicy{engine}
}

func (bp *boardPolicy) CanDelete(ctxUser middleware.CtxUser, id int64) bool {
	return true
}
