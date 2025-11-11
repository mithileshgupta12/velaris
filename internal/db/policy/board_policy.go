package policy

type BoardPolicy interface{}

type boardPolicy struct{}

func NewBoardPolicy() BoardPolicy {
	return &boardPolicy{}
}
