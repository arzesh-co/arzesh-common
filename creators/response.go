package creators

import (
	"github.com/arzesh-co/arzesh-common/errors"
	"strconv"
)

type ResponseOtp struct {
	Meta map[string]any
}

func (otp ResponseOtp) SetPagination(totalCount int64, skip, limit string) {
	if otp.Meta == nil {
		otp.Meta = make(map[string]any)
	}
	current, _ := strconv.ParseInt(skip, 10, 64)
	MLimit, _ := strconv.ParseInt(limit, 10, 64)
	otp.Meta["current"] = current
	if skip == "0" {
		otp.Meta["current"] = 1
		current = 1
	}
	otp.Meta["total"] = totalCount
	if (current * MLimit) < totalCount {
		otp.Meta["next"] = current + 1
	}
	if current > 1 {
		otp.Meta["prev"] = current - 1
	}
}

type apiResponse struct {
	Error   *errors.ResponseErrors `json:"errors,omitempty"`
	Data    any                    `json:"data,omitempty"`
	Message string                 `json:"message,omitempty"`
	Meta    map[string]any         `json:"meta,omitempty"`
}

func New(data any, messageKey string, Error *errors.ResponseErrors, otp *ResponseOtp) *apiResponse {
	res := &apiResponse{}
	//TODO : create message of template
	if messageKey != "" {
		res.Message = messageKey
	}
	if data != nil {
		res.Data = data
	}
	if Error != nil {
		res.Error = Error
	}
	if otp != nil {
		res.Meta = otp.Meta
	}
	return res
}
