package ApiReturn

type ApiReturnCode struct {
	Code int         `json:"code,omitempty"`
	Data interface{} `json:"data,omitempty"`
	Msg  string      `json:"msg,omitempty"`
}

// 构造函数
func ApiReturn(code int, msg string) ApiReturnCode {
	return ApiReturnCode{
		Code: code,
		Msg:  msg,
		//Data: data,
	}
}
func FailureWithMsg(s string) ApiReturnCode {
	failure := Failure
	failure.Msg = s
	return failure
}
func SuccessWithData(data interface{}) ApiReturnCode {
	ok := OK
	ok.Data = data
	return ok
}

var (
	OK                      = ApiReturn(200, "ok")   // 通用成功
	Err                     = ApiReturn(500, "服务异常") // 通用错误
	ErrSystem               = ApiReturn(502, "系统异常") // 系统异常
	Failure                 = ApiReturn(400, "失败")
	SystemBusyness          = ApiReturn(1002, "系统繁忙,请稍后再试")
	ErrCheckParameterFailed = ApiReturn(1003, "参数校验,未通过")
	NoData                  = ApiReturn(1004, "无符合条件数据")
	UnknownErr              = ApiReturn(1005, "未知异常")

	// ErrParam 服务级错误码
	ErrParam            = ApiReturn(10101, "参数错误")
	VerificationNoCount = ApiReturn(10102, "验证码次数已用尽，请一小时后再次")
	ErrVerificationCode = ApiReturn(10103, "验证码验证失败")

	//业务级错误码
)
