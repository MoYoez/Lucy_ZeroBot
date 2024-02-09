package mai

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type Detail struct {
	UserName     string `json:"userName"`
	PlayerRating int    `json:"playerRating"`
	TotalAwake   int    `json:"totalAwake"`
}

type UserLogin struct {
	UserID    int    `json:"userId"`
	ClientID  string `json:"clientId"`
	LoginTime string `json:"loginTime"`
	Detail    Detail `json:"detail"`
}

type ReturnValue struct {
	Data map[string][]UserLogin `json:"returnValue"`
}

func init() {
	engine.OnFullMatchGroup([]string{"bbwj", "bbw几"}).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		if ctx.Event.GroupID != 686575004 && ctx.Event.GroupID != 621400692 {
			return
		}
		getReturnedData, _ := GetSpecifyAttarct(bbwID)
		var returnValue ReturnValue
		err := json.Unmarshal([]byte(getReturnedData), &returnValue)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}
		uniqueMap := make(map[string]bool)
		var getExpectList []int
		for i := range returnValue.Data {
			getInt, err := strconv.Atoi(i)
			if err != nil {
				panic(err)
			}
			getExpectList = append(getExpectList, getInt)
		}
		// sort
		sort.Ints(getExpectList)
		i := 0
		var returnText string
		for range returnValue.Data {
			uniqueList := make([]string, 0, len(uniqueMap))
			for _, userlist := range returnValue.Data[strconv.Itoa(getExpectList[i])] {
				uniqueMap[strconv.Itoa(userlist.UserID)] = true
			}
			for key := range uniqueMap {
				uniqueList = append(uniqueList, key)
			}
			getuserListIDListLength := len(uniqueList)
			getUserFullPC := len(returnValue.Data[strconv.Itoa(getExpectList[i])])
			getReturnText := PlayReturnText(strconv.Itoa(getExpectList[i]), getuserListIDListLength, getUserFullPC)
			returnText = returnText + "\n" + getReturnText
			i = i + 1
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("宝贝王合肥万达广场店的 1 台舞萌 DX：\n"+returnText))
	})
}

func PlayReturnText(timeCost string, playerNumber int, played int) string {
	timeCostToint, _ := strconv.ParseInt(timeCost, 10, 64)
	if timeCostToint > 120 {
		return "\n今天共有 " + strconv.Itoa(playerNumber) + " 位玩家登录了 " + strconv.Itoa(played) + " 次"
	}
	return "在过去的" + timeCost + "分钟内共有 " + strconv.Itoa(playerNumber) + " 玩家登录了 " + strconv.Itoa(played) + " 次"
}
