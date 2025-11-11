package policy

import "github.com/mithileshgupta12/velaris/internal/middleware"

type boardPolicy struct{}

func NewBoardPolicy() Policy {
	return &boardPolicy{}
}

func (bp *boardPolicy) CanDelete(ctxUser middleware.CtxUser, id int64) bool {
	return true
}
