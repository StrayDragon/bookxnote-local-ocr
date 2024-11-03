package handlers

type ErrorResp struct {
	Code    int    `json:"error_code"`           // 错误码
	Message string `json:"error_msg"`            // 错误信息
	XTip    string `json:"error_xtip,omitempty"` // 错误提示（可选）这里作为扩展字段(X*)方便看
}

var (
	ErrRequestLimitResp = ErrorResp{
		Code:    4,
		Message: "Open api request limit reached",
		XTip:    "集群超限额",
	}

	ErrNoPermissionResp = ErrorResp{
		Code:    6,
		Message: "No permission to access data",
		XTip:    "无权限访问该用户数据，请确认已开通相关接口权限",
	}

	ErrIAMAuthResp = ErrorResp{
		Code:    14,
		Message: "IAM Certification failed",
		XTip:    "IAM鉴权失败，请检查签名生成方式是否正确",
	}

	ErrDailyLimitResp = ErrorResp{
		Code:    17,
		Message: "Open api daily request limit reached",
		XTip:    "日调用量超出限制",
	}

	ErrQPSLimitResp = ErrorResp{
		Code:    18,
		Message: "Open api qps request limit reached",
		XTip:    "QPS超出限制",
	}

	ErrTotalLimitResp = ErrorResp{
		Code:    19,
		Message: "Open api total request limit reached",
		XTip:    "总调用量超出限制",
	}

	ErrInvalidParamResp = ErrorResp{
		Code:    100,
		Message: "Invalid parameter",
		XTip:    "请求参数无效",
	}

	ErrInvalidTokenResp = ErrorResp{
		Code:    110,
		Message: "Access token invalid or no longer valid",
		XTip:    "访问令牌无效",
	}

	ErrTokenExpiredResp = ErrorResp{
		Code:    111,
		Message: "Access token expired",
		XTip:    "访问令牌已过期",
	}

	ErrInternalServerResp = ErrorResp{
		Code:    282000,
		Message: "internal error",
		XTip:    "服务器内部错误，请稍后重试",
	}

	ErrEmptyImageResp = ErrorResp{
		Code:    216200,
		Message: "empty image",
		XTip:    "图片数据为空，请检查后重试",
	}

	ErrImageFormatResp = ErrorResp{
		Code:    216201,
		Message: "image format error",
		XTip:    "不支持的图片格式，仅支持PNG、JPG、JPEG、BMP格式",
	}

	ErrImageSizeResp = ErrorResp{
		Code:    216202,
		Message: "image size error",
		XTip:    "图片大小超出限制，base64编码后需小于4M，分辨率不超过4096*4096",
	}
)
