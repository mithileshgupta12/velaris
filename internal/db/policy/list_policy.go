package policy

type ListPolicy interface{}

type listPolicy struct{}

func NewListPolicy() ListPolicy {
	return &listPolicy{}
}
