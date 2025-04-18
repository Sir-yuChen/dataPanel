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

var (
	OK                      = ApiReturn(200, "ok")   // 通用成功
	Err                     = ApiReturn(500, "服务异常") // 通用错误
	ErrSystem               = ApiReturn(502, "系统异常") // 系统异常
	TxErrSystem             = ApiReturn(501, "系统异常") // 事务系统异常
	Failure                 = ApiReturn(400, "失败")
	ErrNoSupport            = ApiReturn(1001, "暂不支持") // 暂时不支持
	SystemBusyness          = ApiReturn(1002, "系统繁忙,请稍后再试")
	ErrCheckParameterFailed = ApiReturn(1003, "参数校验,未通过")
	NoData                  = ApiReturn(1004, "无符合条件数据")
	UnknownErr              = ApiReturn(1005, "未知异常")

	// ErrParam 服务级错误码
	ErrParam            = ApiReturn(10101, "参数有误")
	VerificationNoCount = ApiReturn(10102, "验证码次数已用尽，请一小时后再次")
	ErrVerificationCode = ApiReturn(10103, "验证码验证失败")

	// ErrUserService 模块级错误码 - 用户模块
	NoUserInfo         = ApiReturn(10200, "用户不存在")
	ErrPwd             = ApiReturn(10201, "密码错误")
	ErrCreateToken     = ApiReturn(10202, "token生成异常")
	UserNameExisted    = ApiReturn(10203, "用户名已存在")
	UserPhoneExisted   = ApiReturn(10204, "手机号已被注册")
	LoginExpired       = ApiReturn(10205, "登录已失效,请重新登录")
	UnauthorizedAccess = ApiReturn(10206, "未登录或非法访问")
	DefinedNotType     = ApiReturn(10207, "无法识别对应类型")
	UpdateInfoSame     = ApiReturn(10208, "修改前后不能相同")
	UpdateReLogin      = ApiReturn(10209, "由于修改了登录信息,修改成功,请重新登录")

	//角色权限
	NoPermission       = ApiReturn(10300, "权限不足")
	RoleNameRepeat     = ApiReturn(10301, "角色名不能重复")
	RoleSatatusNo      = ApiReturn(10302, "部分角色不存在/不可用,请刷新重新分配")
	ExistingDepartment = ApiReturn(10303, "当前账户已有归属部门")
	ExistingMenuName   = ApiReturn(10304, "菜单名称存在")

	//消息模板 公告/通讯
	NoGroupInfo          = ApiReturn(10401, "未查询到公告面向人群")
	GroupDisabled        = ApiReturn(10402, "用户群被禁用")
	BulletinFailed       = ApiReturn(10403, "公告新增失败")
	ExistingBulletinName = ApiReturn(10404, "公告名称已经存在")

	//文件
	UploadFailed   = ApiReturn(10450, "文件上传失败,请重试")
	NoFiles        = ApiReturn(10451, "无文件")
	NoRecordFiles  = ApiReturn(10452, "文件记录不存在")
	FileDisabled   = ApiReturn(10453, "文件已禁用/已删除")
	DownloadFailed = ApiReturn(10454, "文件下载失败,请重试")
	NoFaile        = ApiReturn(10455, "文件不存在")
)
