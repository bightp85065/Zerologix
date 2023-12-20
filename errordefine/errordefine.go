package errordefine

// err code  define
const (
	OK                 = 0
	UnknownError       = 1000
	ParameterError     = 1001
	DBError            = 1002
	NotFound           = 1003
	DataError          = 1004
	BasicAccessAuthErr = 2001
	AuthAppKeyErr      = 2002
	AuthAppSecretErr   = 2003
	AuthTokenTypeErr   = 2006
	PermDeniedError    = 4004
)

var ERR_MAP = map[int]string{
	OK:                 "ok",
	UnknownError:       "unknown error",
	ParameterError:     "parameter  error",
	DBError:            "database error",
	NotFound:           "cannot found in database",
	DataError:          "data error",
	BasicAccessAuthErr: "basic access error",
	AuthAppKeyErr:      "appkey is wrong",
	AuthAppSecretErr:   "appsecret is wrong",
	AuthTokenTypeErr:   "auth type is wrong",
	PermDeniedError:    "user has no corresponding access permission",
}
