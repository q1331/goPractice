package main

type argError struct{
	msg string
}

func (e argError) Error() string{
	return e.msg
}