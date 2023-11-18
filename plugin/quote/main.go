package quote

import (
	"fmt"
	"hash/crc64"
	"regexp"
	"strconv"
	"time"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

var engine = control.Register("quote", &ctrl.Options[*zero.Ctx]{
	DisableOnDefault:  false,
	Help:              "Make A Quote! 记录群友发言",
	PrivateDataFolder: "quote",
})

func init() {
	engine.OnRegex(`^(\[CQ:reply,id=(.*)\])\s/quote$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getPatternUserMessageID := ctx.State["regex_matched"].([]string)[2]
		rsp := ctx.CallAction("get_msg", zero.Params{
			"message_id": getPatternUserMessageID,
		}).Data.String()
		SenderUserID := gjson.Get(rsp, "sender.user_id").Int()
		MessageRaw := gjson.Get(rsp, "message").String()
		formatMessageRaw := message.UnescapeCQCodeText(MessageRaw)
		InsertUserCorruptionTarget(Corruption{
			TrackID:   TrackIDGenerator(ctx.Event.GroupID),
			QQ:        SenderUserID,
			Msg:       formatMessageRaw,
			HandlerQQ: ctx.Event.UserID,
			Time:      time.Now().Unix(),
		}, ctx.Event.GroupID)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("成功~"))
	})
	engine.OnFullMatch("/quote").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		// make a quote
		getRandGeneratedData := IndexOfDataCorruption(ctx.Event.GroupID)
		if getRandGeneratedData.TrackID == 0 {
			ctx.Send(message.Text("找不到Quote，或许是本群没有x"))
			return
		}
		ctx.Send(fmt.Sprintf("TrackerID: %d\nTrackerUser: %s(%d)\nTime: %s\n%s(%d): \n%s", getRandGeneratedData.TrackID, ctx.CardOrNickName(getRandGeneratedData.HandlerQQ), getRandGeneratedData.HandlerQQ, time.Unix(getRandGeneratedData.Time, 0).Format("2006-01-02 15:04:05"), ctx.CardOrNickName(getRandGeneratedData.QQ), getRandGeneratedData.QQ, getRandGeneratedData.Msg))
	})
	engine.OnRegex(`^(\[CQ:reply,id=(.*)\])\s/quote\sremove$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getPatternUserMessageID := ctx.State["regex_matched"].([]string)[2]
		rsp := ctx.CallAction("get_msg", zero.Params{
			"message_id": getPatternUserMessageID,
		}).Data.String()
		// get msg string.
		MessageRaw := gjson.Get(rsp, "message").String()
		formatMessageRaw := message.UnescapeCQCodeText(MessageRaw)
		// use pattern to get track id.
		trackerIDPattern, err := regexp.Compile(`^TrackerID:\s(\d*)`)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("找不到对应quote，可能是该Quote已经被删除("))
			return
		}
		getString := trackerIDPattern.FindStringSubmatch(formatMessageRaw)[1]
		fmt.Printf(getString)
		if getString == "" {
			ctx.Send("请回复对应需要删除的quote来进行此操作")
			return
		}
		// remove string must in person or admin.
		getInfo := ReferDataCorruption(ctx.Event.GroupID, getString)
		if ctx.Event.UserID != getInfo.QQ && ctx.Event.UserID != getInfo.HandlerQQ && !zero.AdminPermission(ctx) {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("大概你不是管理员 | quote制作者/生成者 此操作无效"))
			return
		}
		// remove it
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("删除完成~"))
		err = RemoveIndexData(getString, ctx.Event.GroupID)
		if err != nil {
			panic(err)
		}
	})
}

func TrackIDGenerator(groupID int64) int64 {
	return int64(crc64.Checksum(helper.StringToBytes(strconv.FormatInt(groupID+time.Now().Unix(), 10)), crc64.MakeTable(crc64.ISO)))
}
