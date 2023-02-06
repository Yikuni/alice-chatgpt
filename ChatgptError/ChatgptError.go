package ChatgptError

type ChatgptError struct {
	msg string
	error
}

func Err(msg string) ChatgptError {
	return ChatgptError{msg: msg}
}
func (err ChatgptError) Error() string {
	return err.msg
}
