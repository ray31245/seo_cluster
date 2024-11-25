package model

const (
	// context
	ContextEdit  ApiContext = "edit"
	ContextView  ApiContext = "view"
	ContextEmbed ApiContext = "embed"

	// order
	OrderAsc  ApiOrder = "asc"
	OrderDesc ApiOrder = "desc"
)

type (
	ApiContext string
	ApiOrder   string
)

type BasicAuthentication struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	IsAnonymous bool   `json:"isAnonymous"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Status int `json:"status"`
	} `json:"data"`
}
