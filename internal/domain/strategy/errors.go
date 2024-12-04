package strategy

type PackagingError struct {
	Msg string
}

func (e PackagingError) Error() string {
	return e.Msg
}
