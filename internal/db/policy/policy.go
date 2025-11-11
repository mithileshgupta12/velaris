package policy

import "github.com/mithileshgupta12/velaris/internal/middleware"

type Policy interface {
	CanDelete(ctxUser middleware.CtxUser, id int64) bool
}
