package utils

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func APIResponseError(message string, data interface{}) APIResponse {
	return APIResponse{
		Status:  "error",
		Message: message,
		Data:    data,
	}
}

func APIResponseSuccess(message string, data interface{}) APIResponse {
	return APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}
