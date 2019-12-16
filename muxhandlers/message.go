package muxhandlers

import (
	"encoding/json"
	"strconv"

	"github.com/Mtbcooler/outrun/db"
	"github.com/Mtbcooler/outrun/emess"
	"github.com/Mtbcooler/outrun/enums"
	"github.com/Mtbcooler/outrun/helper"
	"github.com/Mtbcooler/outrun/obj"
	"github.com/Mtbcooler/outrun/requests"
	"github.com/Mtbcooler/outrun/responses"
	"github.com/Mtbcooler/outrun/status"
)

func GetMessageList(helper *helper.Helper) {
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("error getting calling player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	if player.Messages == nil {
		player.Messages = []obj.Message{}
	}
	if player.OperatorMessages == nil {
		player.OperatorMessages = []obj.OperatorMessage{}
	}
	db.SavePlayer(player)
	// response := responses.DefaultMessageList(baseInfo)
	response := responses.MessageList(baseInfo, player.Messages, player.OperatorMessages)
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
}

func GetMessage(helper *helper.Helper) {
	data := helper.GetGameRequest()
	var request requests.GetMessageRequest
	err := json.Unmarshal(data, &request)
	if err != nil {
		helper.InternalErr("Error unmarshalling", err)
		return
	}
	player, err := helper.GetCallingPlayer()
	if err != nil {
		helper.InternalErr("error getting calling player", err)
		return
	}
	baseInfo := helper.BaseInfo(emess.OK, status.OK)
	if player.Messages == nil {
		player.Messages = []obj.Message{}
	}
	if player.OperatorMessages == nil {
		player.OperatorMessages = []obj.OperatorMessage{}
	}

	presentList := []obj.Present{}

	switch messageIds := request.MessageIDs.(type) {
	case []interface{}:
		helper.DebugOut("%v", messageIds)
		for _, msgid := range messageIds {
			helper.DebugOut("Accepting message ID %v", msgid)
			present := player.AcceptMessage(int64(msgid.(float64))) // TODO: why does Go think this is a float64 and not an int64?
			if present != nil {
				presentList = append(presentList, present.(obj.Present))
			}
		}
	case string:
		helper.DebugOut("No messages to accept")
	default:
		helper.Warn("Unexpected type of request.MessageIDs")
	}

	switch operatorMessageIds := request.OperatorMessageIDs.(type) {
	case []interface{}:
		helper.DebugOut("%v", operatorMessageIds)
		player.CleanUpExpiredOperatorMessages()
		for _, omsgid := range operatorMessageIds {
			helper.DebugOut("Accepting operator message ID %v", omsgid)
			present := player.AcceptOperatorMessage(int64(omsgid.(float64))) // TODO: why does Go think this is a float64 and not an int64?
			if present != nil {
				presentList = append(presentList, present.(obj.Present))
			}
		}
	case string:
		helper.DebugOut("No operator messages to accept")
	default:
		helper.Warn("Unexpected type of request.OperatorMessageIDs")
	}

	helper.DebugOut("%v", presentList)
	for _, currentPresent := range presentList {
		itemid := strconv.Itoa(int(currentPresent.ItemID))
		helper.DebugOut("Present: %s", itemid)
		helper.DebugOut("Present amount: %v", currentPresent.NumItem)
		if itemid[:2] == "12" { // ID is an item
			// check if the item is already in the player's inventory
			for _, item := range player.PlayerState.Items {
				if item.ID == itemid { // item found, increment amount
					item.Amount += currentPresent.NumItem
					break
				}
			}
		} else if itemid == strconv.Itoa(enums.ItemIDRing) { // Rings
			player.PlayerState.NumRings += currentPresent.NumItem
		} else if itemid == strconv.Itoa(enums.ItemIDRedRing) { // Red rings
			player.PlayerState.NumRedRings += currentPresent.NumItem
		} else if itemid == strconv.Itoa(enums.ItemIDEnergy) { // Revive tokens
			player.PlayerState.Energy += currentPresent.NumItem
		} else if itemid == strconv.Itoa(enums.IDSpecialEgg) {
			player.PlayerState.ChaoEggs += currentPresent.NumItem
		} else if itemid == strconv.Itoa(enums.IDRouletteTicketPremium) {
			player.PlayerState.NumChaoRouletteTicket += currentPresent.NumItem
		} else if itemid == strconv.Itoa(enums.IDRouletteTicketItem) {
			player.PlayerState.NumRouletteTicket += currentPresent.NumItem
		} else if itemid[:2] == "40" { // ID is a Chao
			chaoIndex := player.IndexOfChao(itemid)
			if chaoIndex == -1 { // chao index not found, should never happen
				helper.InternalErr("cannot get index of chao '"+strconv.Itoa(chaoIndex)+"'", err)
				return
			}
			if player.ChaoState[chaoIndex].Status == enums.ChaoStatusNotOwned {
				// earn the Chao
				player.ChaoState[chaoIndex].Status = enums.ChaoStatusOwned
				player.ChaoState[chaoIndex].Acquired = 1
				player.ChaoState[chaoIndex].Level = 0
			}
			player.ChaoState[chaoIndex].Level += currentPresent.NumItem
			if player.ChaoState[chaoIndex].Level > 10 { // if max chao level
				player.ChaoState[chaoIndex].Level = 10                        // reset to maximum
				player.ChaoState[chaoIndex].Status = enums.ChaoStatusMaxLevel // set status to MaxLevel
			}
		} else {
			helper.Out("Unknown present ID %s", itemid)
		}
	}
	var response interface{}
	if baseInfo.StatusCode == status.OK {
		response = responses.GetMessage(baseInfo, player, presentList, player.GetAllOperatorMessageIDs(), player.GetAllOperatorMessageIDs())
	} else {
		response = responses.NewBaseResponse(baseInfo)
	}
	err = helper.SendResponse(response)
	if err != nil {
		helper.InternalErr("Error sending response", err)
	}
	db.SavePlayer(player)
}
