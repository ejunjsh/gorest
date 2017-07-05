package gorest


type NoFoundError struct {}

func (e NoFoundError) Error() string  {
	return "404 No Found"
}


type InternalError struct {
	Err error
	Message string
}


func (e InternalError) Error() string  {
	if e.Err!=nil{
		return e.Err.Error()
	}
	return e.Message
}