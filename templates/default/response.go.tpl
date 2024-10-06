package utils

type JSONResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(message string, data interface{}, code int) JSONResponse {
	return JSONResponse{
		Status:  "success",
		Message: message,
		Code:    code,
		Data:    data,
	}
}

func Error(message string, code int) JSONResponse {
	return JSONResponse{
		Status:  "error",
		Message: message,
		Code:    code,
	}
}