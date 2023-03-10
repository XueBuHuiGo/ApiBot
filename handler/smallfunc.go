package handler

import (
	"ApiBot/util"
	"fmt"
	"regexp"
)

// Jiujiu 啾啾
func Jiujiu(userId int64) (string, bool, error) {
	rsg := fmt.Sprintf("[CQ:image,file=http://ovooa.com/API/face_kiss/?QQ=%d,cache=0]", userId)
	return rsg, false, nil
}

// Yiyan 一言
func Yiyan() (string, bool, error) {
	rsp, err := util.MyRequest("GET", "https://v.api.aa1.cn/api/yiyan/index.php", nil, "")
	if err != nil {
		return "一言 请求失败", true, err
	}
	re := regexp.MustCompile("<p>(.*)</p>")
	p := re.FindStringSubmatch(string(rsp))
	rsg := p[1]
	return rsg, true, nil
}

// Dogdiary 舔狗日记
func Dogdiary() (string, bool, error) {
	rsp, err := util.MyRequest("GET", "https://v.api.aa1.cn/api/tiangou/index.php", nil, "")
	if err != nil {
		return "舔狗日记 请求失败", true, err
	}
	re := regexp.MustCompile("<p>(.*)</p>")
	p := re.FindStringSubmatch(string(rsp))
	rsg := p[1]
	return rsg, true, nil
}
