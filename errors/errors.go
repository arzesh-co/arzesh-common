package errors

type ResponseErrors struct {
	ErrorType struct {
		Title string `json:"title"`
		Desc  string `json:"desc"`
		Url   string `json:"url"`
	} `json:"error_type"`
	StatusKey string `json:"status_key"`
	Detail    string `json:"detail"`
	Title     string `json:"title"`
	HelpUrl   string `json:"help_url"`
}

func New(key, entity, params map[string]string) *ResponseErrors {

}
