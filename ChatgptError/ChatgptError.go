package ChatgptError

type ChatgptError struct {
	msg string
	error
}

type ExceededQuotaException struct {
	ChatgptError
}

func Err(msg string) error {
	if msg == "You exceeded your current quota,please check your plan and billing details. " {
		return ExceededQuotaException{ChatgptError{msg: msg}}
	}
	return ChatgptError{msg: msg}
}
func (err ChatgptError) Error() string {
	return err.msg
}
