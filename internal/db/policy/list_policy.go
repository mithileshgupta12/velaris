package policy

import (
	"github.com/mithileshgupta12/velaris/internal/middleware"
	"xorm.io/xorm"
)

type listPolicy struct {
	engine *xorm.Engine
}

func NewListPolicy(engine *xorm.Engine) Policy {
	return &listPolicy{engine}
}

func (lp *listPolicy) userOwnsList(ctxUser middleware.CtxUser, id int64) (bool, error) {
	exists, err := lp.engine.
		Alias("l").
		Where("l.id = ?", id).
		Join("INNER", "boards b", "b.id = l.board_id").
		Where("b.user_id = ?", ctxUser.ID).
		Exist()
	if err != nil {
		return false, err
	}
	if !exists {
		return false, nil
	}

	return true, nil
}

func (lp *listPolicy) CanView(ctxUser middleware.CtxUser, id int64) (bool, error) {
	return lp.userOwnsList(ctxUser, id)
}

func (lp *listPolicy) CanCreate(ctxUser middleware.CtxUser, id int64) (bool, error) {
	return true, nil
}

func (lp *listPolicy) CanUpdate(ctxUser middleware.CtxUser, id int64) (bool, error) {
	return lp.userOwnsList(ctxUser, id)
}

func (lp *listPolicy) CanDelete(ctxUser middleware.CtxUser, id int64) (bool, error) {
	return lp.userOwnsList(ctxUser, id)
}
