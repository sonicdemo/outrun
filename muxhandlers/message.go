package muxhandlers

import (
    "github.com/Mtbcooler/outrun/emess"
    "github.com/Mtbcooler/outrun/helper"
    "github.com/Mtbcooler/outrun/responses"
    "github.com/Mtbcooler/outrun/status"
)

func GetMessageList(helper *helper.Helper) {
    // player agnostic
    baseInfo := helper.BaseInfo(emess.OK, status.OK)
    response := responses.DefaultMessageList(baseInfo)
    err := helper.SendResponse(response)
    if err != nil {
        helper.InternalErr("Error sending response", err)
    }
}
