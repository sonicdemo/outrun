package muxhandlers

import (
	"encoding/json"

	"github.com/RunnersRevival/outrun/emess"
	"github.com/RunnersRevival/outrun/helper"
	"github.com/RunnersRevival/outrun/requests"
	"github.com/RunnersRevival/outrun/responses"
	"github.com/RunnersRevival/outrun/status"
)

func GetWeeklyLeaderboardOptions(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.LeaderboardRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}
	mode := request.Mode
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultWeeklyLeaderboardOptions(baseInfo, mode)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetWeeklyLeaderboardEntries(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.LeaderboardEntriesRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("Error getting calling player", err)
		return
	}
	mode := request.Mode
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultWeeklyLeaderboardEntries(baseInfo, player, mode)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetLeagueData(helper *helper.Helper) {
	recv := helper.GetGameRequest()
	var request requests.LeaderboardRequest
	err := json.Unmarshal(recv, &request)
	if err != nil {
		helper.Err("Error unmarshalling", err)
		return
	}
	mode := request.Mode
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	response := responses.DefaultLeagueData(baseInfo, mode)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}
