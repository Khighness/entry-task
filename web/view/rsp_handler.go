package view

import (
	"encoding/json"
	"github.com/Khighness/entry-task/web/common"
	"net/http"
)

// @Author Chen Zikang
// @Email  zikang.chen@shopee.com
// @Since  2022-02-23

// HandleBizSuccess 处理业务成功结果
func HandleBizSuccess(w http.ResponseWriter, data interface{}) {
	response, _ := json.Marshal(common.HttpResponse{
		Code:    common.HttpSuccessCode,
		Message: common.HttpSuccessMessage,
		Data:    data,
	})
	_, _ = w.Write(response)
}

// HandleBizError 处理业务错误结果
func HandleBizError(w http.ResponseWriter, data interface{}) {
	response, _ := json.Marshal(common.HttpResponse{
		Code:    common.HttpErrorCode,
		Message: common.HttpErrorMessage,
		Data:    data,
	})
	_, _ = w.Write(response)
}

// HandleErrorServerBusy 处理服务繁忙
func HandleErrorServerBusy(w http.ResponseWriter) {
	response, _ := json.Marshal(common.HttpResponse{
		Code:    common.HttpErrorServerBusyCode,
		Message: common.HttpErrorServerBusyMessage,
	})
	_, _ = w.Write(response)
}

// HandleErrorRpcRequest 处理RPC请求错误
func HandleErrorRpcRequest(w http.ResponseWriter) {
	response, _ := json.Marshal(common.HttpResponse{
		Code:    common.HttpErrorRpcRequestCode,
		Message: common.HttpErrorRpcRequestMessage,
	})
	_, _ = w.Write(response)
}

// HandleErrorRpcResponse 处理RPC结果错误
func HandleErrorRpcResponse(w http.ResponseWriter, code int32, msg string) {
	response, _ := json.Marshal(common.HttpResponse{
		Code:    code,
		Message: msg,
	})
	_, _ = w.Write(response)
}

// HandleMethodError 处理方法错误
func HandleMethodError(w http.ResponseWriter, data interface{}) {
	response, _ := json.Marshal(common.HttpResponse{
		Code:    http.StatusMethodNotAllowed,
		Message: "Method Not Allowed",
		Data:    data,
	})
	_, _ = w.Write(response)
}

// HandleRequestError 处理请求错误
func HandleRequestError(w http.ResponseWriter, data interface{}) {
	response, _ := json.Marshal(common.HttpResponse{
		Code:    http.StatusBadRequest,
		Message: "Bad Request",
		Data:    data,
	})
	_, _ = w.Write(response)
}
