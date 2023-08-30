package constant

const (
	DATABASE_NAME = "public"
	SCORE_TABLE   = "mytable"
	USER_TABLE    = "user"
	RSAKeyTimeout = 1800 // 单位 秒
	RSAKeyBitsize = 2048
)

// 错误码
const (
	SUCCESS               = 0
	UNKNOWN               = -1
	AUTH_PARSE_ERR        = -1000
	AUTH_KEY_TIMEOUT      = -1001
	AUTH_RSAKEY_NOT_FOUND = -1002
)
