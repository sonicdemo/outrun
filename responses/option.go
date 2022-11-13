package responses

import (
	"github.com/RunnersRevival/outrun/netobj"
	"github.com/RunnersRevival/outrun/responses/responseobjs"
)

type OptionUserResultResponse struct {
	BaseResponse
	OptionUserResult netobj.OptionUserResult `json:"optionUserResult"`
}

func OptionUserResult(base responseobjs.BaseInfo, optionUserResult netobj.OptionUserResult) OptionUserResultResponse {
	baseResponse := NewBaseResponse(base)
	out := OptionUserResultResponse{
		baseResponse,
		optionUserResult,
	}
	return out
}
