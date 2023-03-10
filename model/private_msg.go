package model

import (
	"ApiBot/handler"
)

func PrivateMsg(Data map[string]interface{}) (string, bool, interface{}) {
	message := Data["NewMessage"].(string)
	var (
		rsg string
		//AutoEscape bool
		err interface{}
	)
	rsg, err = handler.ChatGPT(Data, message)
	if err != nil {
		return "GPT 请求失败", false, err
	}
	return rsg, false, nil
}
