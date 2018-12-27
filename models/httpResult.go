package models

type HttpResult struct {
	ResultCode string
	ResultMsg string
	ResultData interface{}
}

func init() {
}

func GetJsonResult(Data interface{}) (h *HttpResult) {
	r := HttpResult{"200", "成功", Data}
	return &r
}

func GetErrorResult(code string,msg string) (h *HttpResult) {
	r := HttpResult{code, msg, ""}
	return &r
}