package domain

type (
	Response struct {
		Data interface{} `json:"data,omitempty"`
		Meta Meta        `json:"meta,omitempty"`
	}

	Meta struct {
		Code    string `json:"code,omitempty"`
		Message string `json:"message,omitempty"`
	}
)

type (
	PaginationResponse struct {
		Number        int    `json:"number,omitempty"`
		Size          int    `json:"size,omitempty"`
		TotalElements int    `json:"total_elements,omitempty"`
		TotalPages    int    `json:"total_pages,omitempty"`
		Sort          string `json:"sort,omitempty"`
		SortBy        string `json:"sort_by,omitempty"`
	}

	ValidationResponse struct {
		Result bool `json:"result"`
	}
)

type (
	DeleteRequest struct {
		Id           uint   `json:"id"`
		LoggedUserId uint   `json:"logged_user_id"`
		LoggedUser   string `json:"logged_user"`
	}

	FindByIdRequest struct {
		Id uint `json:"id"`
	}

	SearchRequest struct {
		TextSearch string `json:"text_search"`
		Start      uint   `json:"start"`
		Limit      uint   `json:"limit"`
		SortBy     string `json:"sort_by"`
		Sort       string `json:"sort"`
	}

	SessionRequest struct {
		UserId       uint   `json:"user_id" swaggerignore:"true"`
		Token        string `json:"token" swaggerignore:"true"`
		RefreshToken string `json:"refresh_token" swaggerignore:"true"`
		Username     string `json:"username" swaggerignore:"true"`
		ChannelId    string `json:"channel_id" swaggerignore:"true"`
		ApiKey       string `json:"api_key" swaggerignore:"true"`
		AccessToken  string `json:"access_token" swaggerignore:"true"`
	}
)

func BusinessErrorResponse(e *BussinessError) *Response {
	return &Response{
		Meta: Meta{Code: e.ErrorCode, Message: e.ErrorMessage},
		Data: nil,
	}
}

func TechnicalErrorResponse(e *TechnicalError) *Response {
	return &Response{
		Meta: Meta{Code: GeneralError, Message: SomethingWrong},
		Data: e,
	}
}

func DefaultSuccessResponse(m string, data interface{}) *Response {
	return &Response{
		Meta: Meta{Code: SuccessCode, Message: m},
		Data: data,
	}
}

func CustomResponse(c string, m string, data interface{}) *Response {
	return &Response{
		Meta: Meta{Code: c, Message: m},
		Data: data,
	}
}
