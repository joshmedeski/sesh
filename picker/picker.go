package picker

type Picker interface {
	Pick()
}

type RealPicker struct {
}

func NewPicker() *RealPicker {
	return &RealPicker{
	}
}
