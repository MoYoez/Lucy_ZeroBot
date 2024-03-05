package funwork

import (
	"math/rand"
	"time"

	"github.com/FloatTech/zbputils/ctxext"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

// debug
var fail = "获取精华消息失败喵~可能是这条信息在数据库中无法查询~"

func init() {
	engine.OnFullMatch("今日龙王", zero.OnlyGroup).Limit(ctxext.LimitByGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		list := ctx.GetGroupHonorInfo(ctx.Event.GroupID, "talkative")
		temp := list.String()
		id := gjson.Get(temp, "current_talkative.user_id")
		name := ctx.CardOrNickName(id.Int())
		if name == "" {
			ctx.SendChain(message.Text("今日没有龙王～"))
			return
		}
		ctx.SendChain(message.Text("今日的龙王是~", name, "哦"))
	})
	
	engine.OnFullMatch("随机本群精华消息").SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		essenceList := ctx.GetThisGroupEssenceMessageList()
		essenceCount := len(essenceList.Array())
		if essenceCount == 0 {
			ctx.Send(fail)
		} else {
			IDx := rand.Intn(essenceCount)
			essenceMessage := essenceList.Array()[IDx]
			var (
				nickname = gjson.Get(essenceMessage.Raw, "sender_nick")
				msID     = gjson.Get(essenceMessage.Raw, "message_id")
			)
			ctx.GetGroupMessageHistory(ctx.Event.GroupID, msID.Int())
			ms := ctx.GetMessage(message.NewMessageIDFromInteger(msID.Int()))
			reportText := message.Text("Lucy抓到了这一条消息~\nUsername: ", nickname)
			report := make(message.Message, len(ms.Elements))
			report = append(report, reportText)
			report = append(report, ms.Elements...)
			deleteme := ctx.Send(report)
			time.Sleep(time.Second * 40)
			ctx.DeleteMessage(deleteme)
		}
	})
}
