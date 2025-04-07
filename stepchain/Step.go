package stepchain

type Step interface {
	GetName() string
	Update() (Step, error)
}
