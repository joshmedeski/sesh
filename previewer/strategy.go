package previewer

type PreviewStrategy interface {
	Execute(name string) (string, error)
}
