package creators

type ResponseOtp struct {
	Meta map[string]any
}

func (otp ResponseOtp) SetPagination(totalCount, skip, limit string) {

}

func New(data any, messageKey, ErrorKey string, otp ...map[string]any) {

}
