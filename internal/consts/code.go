package consts

import "github.com/gogf/gf/v2/errors/gcode"

const (
	CodeOK            = 0
	CodeInvalidParams = 1001
	CodeNotFound      = 1002
	CodeInternalError = 1003
	CodeUserExists    = 1004
	CodeBatchTooLarge = 1005
)

var (
	OK            = gcode.New(CodeOK, "OK", nil)
	InvalidParams = gcode.New(CodeInvalidParams, "Invalid parameters", nil)
	NotFound      = gcode.New(CodeNotFound, "Resource not found", nil)
	InternalError = gcode.New(CodeInternalError, "Internal server error", nil)
	UserExists    = gcode.New(CodeUserExists, "User already exists", nil)
	BatchTooLarge = gcode.New(CodeBatchTooLarge, "Batch size too large", nil)
)
