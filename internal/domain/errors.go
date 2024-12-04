package domain

type OrderError struct {
	Msg string
}

func (e OrderError) Error() string {
	return e.Msg
}
