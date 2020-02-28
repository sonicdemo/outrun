package responses

import (
	"github.com/Mtbcooler/outrun/consts"
	"github.com/Mtbcooler/outrun/responses/responseobjs"
)

type BaseResponse struct {
	responseobjs.BaseInfo
	AssetsVersion     string `json:"assets_version"`
	ClientDataVersion string `json:"client_data_version"`
	DataVersion       string `json:"data_version"`
	InfoVersion       string `json:"info_version"`
	Version           string `json:"version"`
}

func NewBaseResponse(base responseobjs.BaseInfo) BaseResponse {
	return BaseResponse{
		base,
		"049",
		"2.0.3",
		"15",
		"017",
		"2.0.3",
	}
}

func NewBaseResponseV(base responseobjs.BaseInfo, gameVersion string) BaseResponse {
	return BaseResponse{
		base,
		consts.DataVersionForGameVersion[gameVersion],
		gameVersion,
		"15",
		"017",
		gameVersion,
	}
}
