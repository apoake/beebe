package model

var SUCCESS *ErrorCode = &ErrorCode{Code: 200, Msg: "SUCCESS"}
var SYSTEM_ERROR *ErrorCode = &ErrorCode{Code: 500, Msg: "SYSTEM_ERROR"}
var PARAMETER_INVALID *ErrorCode = &ErrorCode{Code: 600, Msg: "PARAMETER_INVALID"}
var USERNAME_PASSWORD_ERROR *ErrorCode = &ErrorCode{Code: 601, Msg: "USERNAME_PASSWORD_ERROR"}
var USER_NO_LOGIN *ErrorCode = &ErrorCode{Code: 602, Msg: "USER_NO_LOGIN"}
var USER_ALREADY_LOGIN *ErrorCode = &ErrorCode{Code: 603, Msg: "USER_ALREADY_LOGIN"}
var USER_REGISTER_ERROR *ErrorCode = &ErrorCode{Code: 604, Msg: "USER_REGISTER_ERROR"}
var USER_ACCOUNT_EXIST *ErrorCode = &ErrorCode{Code: 604, Msg: "USER_ACCOUNT_EXIST"}

type ErrorCode struct {
	Code int
	Msg  string
}