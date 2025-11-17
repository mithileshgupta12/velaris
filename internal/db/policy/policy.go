package policy

import (
	"github.com/mithileshgupta12/velaris/internal/middleware"
	"xorm.io/xorm"
)

// Policy defines the authorization interface for checking user permissions
// on various CRUD operations for a resource.
type Policy interface {
	// CanView checks if a user can view a resource with the given id.
	// Returns true if the user has view permission, false otherwise.
	CanView(ctxUser middleware.CtxUser, id int64) (bool, error)

	// CanCreate checks if a user can create a resource with the given id.
	// Returns true if the user has create permission, false otherwise.
	CanCreate(ctxUser middleware.CtxUser, id int64) (bool, error)

	// CanUpdate checks if a user can update a resource with the given id.
	// Returns true if the user has update permission, false otherwise.
	CanUpdate(ctxUser middleware.CtxUser, id int64) (bool, error)

	// CanDelete checks if a user can delete a resource with the given id.
	// Returns true if the user has delete permission, false otherwise.
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
