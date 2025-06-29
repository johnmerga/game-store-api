package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func JSON(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func Success(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

func SuccessWithMessage(data interface{}, message string) Response {
	return Response{
		Success: true,
		Data:    data,
		Message: message,
	}
}

func Error(err interface{}) Response {
	return Response{
		Success: false,
		Error:   err,
	}
}

func ErrorWithMessage(err interface{}, message string) Response {
	return Response{
		Success: false,
		Error:   err,
		Message: message,
	}
}

func Paginated(data interface{}, page, limit, total int) PaginatedResponse {
	totalPages := (total + limit - 1) / limit

	return PaginatedResponse{
		Success: true,
		Data:    data,
		Pagination: Pagination{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}
}
