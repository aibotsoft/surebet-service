package handler

type SurebetError struct {
	cause     error
	msg       string
	permanent bool
}

func (s SurebetError) Error() string {
	if s.cause != nil {
		return s.msg + ": " + s.cause.Error()
	}
	return s.msg
}
