package gpterror

type ChatgptError struct {
	Msg string
	error
}

type ExceededQuotaException struct {
	ChatgptError
}

func Err(msg string) error {
	if msg == "You exceeded your current quota, please check your plan and billing details." {
		return ExceededQuotaException{ChatgptError{Msg: msg}}
	}
	return ChatgptError{Msg: msg}
}
func (err ChatgptError) Error() string {
	return err.Msg
}
