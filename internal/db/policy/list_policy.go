package policy

import "github.com/mithileshgupta12/velaris/internal/middleware"

type listPolicy struct{}

func NewListPolicy() Policy {
	return &listPolicy{}
}

func (lp *listPolicy) CanView(ctxUser middleware.CtxUser, id int64) (bool, error) {
	return true, nil
}

func (lp *listPolicy) CanCreate(ctxUser middleware.CtxUser, id int64) (bool, error) {
	return true, nil
}

func (lp *listPolicy) CanUpdate(ctxUser middleware.CtxUser, id int64) (bool, error) {
	return true, nil
}

func (lp *listPolicy) CanDelete(ctxUser middleware.CtxUser, id int64) (bool, error) {
	return true, nil
}
