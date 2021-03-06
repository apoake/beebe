package model

/**
	common error
 */
var SUCCESS *ErrorCode = &ErrorCode{Code: 200, Msg: "SUCCESS"}
var SYSTEM_ERROR *ErrorCode = &ErrorCode{Code: 500, Msg: "SYSTEM_ERROR"}
var PARAMETER_INVALID *ErrorCode = &ErrorCode{Code: 600, Msg: "PARAMETER_INVALID"}
/**
	user error
 */
var USERNAME_PASSWORD_ERROR *ErrorCode = &ErrorCode{Code: 601, Msg: "USERNAME_PASSWORD_ERROR"}
var USER_NO_LOGIN *ErrorCode = &ErrorCode{Code: 602, Msg: "USER_NO_LOGIN"}
var USER_ALREADY_LOGIN *ErrorCode = &ErrorCode{Code: 603, Msg: "USER_ALREADY_LOGIN"}
var USER_REGISTER_ERROR *ErrorCode = &ErrorCode{Code: 604, Msg: "USER_REGISTER_ERROR"}
var USER_ACCOUNT_EXIST *ErrorCode = &ErrorCode{Code: 605, Msg: "USER_ACCOUNT_EXIST"}
var USER_NO_RIGHT *ErrorCode = &ErrorCode{Code: 606, Msg: "USER_ACCOUNT_EXIST"}
var USER_PASSWORD_DISAGREE *ErrorCode = &ErrorCode{Code: 607, Msg: "USER_PASSWORD_DISAGREE"}
var TEAM_NO_EXIST *ErrorCode = &ErrorCode{Code: 608, Msg: "TEAM_NO_EXIST"}
var TEAM_ALREADY_JOIN *ErrorCode = &ErrorCode{Code: 609, Msg: "TEAM_ALREADY_JOIN"}
/**
	project error
 */
var PROJECT_CREATE_ERROR *ErrorCode = &ErrorCode{Code: 701, Msg: "PROJECT_CREATE_ERROR"}
var PROJECT_ACTION_NOT_FIND *ErrorCode = &ErrorCode{Code: 702, Msg: "PROJECT_ACTION_NOT_FIND"}
var WORKSPACE_ALREADY_IN *ErrorCode = &ErrorCode{Code: 703, Msg: "WORKSPACE_ALREADY_IN"}


/**
	team error
 */
var TEAM_NOT_SELF *ErrorCode = &ErrorCode{Code: 801, Msg: "TEAM_NOT_SELF"}
var TEAM_NOT_EXIST *ErrorCode = &ErrorCode{Code: 802, Msg: "TEAM_NOT_EXIST"}


type ErrorCode struct {
	Code int
	Msg  string
}