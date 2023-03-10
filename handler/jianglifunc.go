package handler

import (
	"ApiBot/util"
	"fmt"
	"github.com/goccy/go-json"
	"math/rand"
	"time"
)

func Cos() (string, bool, error) {
	rsg := "[CQ:image,file=http://43.139.139.17:778/P/cos/index.php,cache=0]"
	return rsg, false, nil
}

func Kanmn() (string, bool, error) {
	rsg := "[CQ:image,file=https://v.api.aa1.cn/api/pc-girl_bz/index.php?wpon=ro38d57y8rhuwur3788y3rd,cache=0]"
	return rsg, false, nil
}

func Xjj() (string, bool, error) {
	source := rand.NewSource(time.Now().Unix())
	random := rand.New(source)
	Xjjapi := [...]string{"https://zj.v.api.aa1.cn/api/video_dyv2"}
	rsp, _ := util.MyRequest("GET", Xjjapi[random.Intn(len(Xjjapi))], nil, "")
	var rspmap map[string]interface{}
	if err := json.Unmarshal(rsp, &rspmap); err != nil {
		return "小姐姐 请求失败", true, err
	}
	rsg := fmt.Sprintf("[CQ:video,file=%v]", rspmap["url"].(string))
	return rsg, false, nil
}
