package common

type BusinessCode int
type Result struct {
	Code BusinessCode `json:"code"`
	Data any `json:"data"`
	Msg string `json:"msg"`
}


func (r *Result) Success(data interface{}, msg ...string) *Result {
	r.Code = 200 
	if len(msg) > 0 {
		r.Msg = msg[0]
	} else {
		r.Msg = "success"
	}
	r.Data = data
	return r
}

func (r *Result) Fail(code BusinessCode, msg string) *Result {
	r.Code = code 
	r.Msg = msg
	return r
}