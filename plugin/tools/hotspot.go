// Package tools for tools 参考了 Cha0sIDL 的 zbp hotspot热点
package tools

import (
	"strconv"

	"github.com/FloatTech/AnimeAPI/bilibili"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	engine.OnFullMatch("微博热搜", zero.OnlyToMe).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		rsp := "微博实时热榜:\n"
		url := "https://api.weibo.cn/2/guest/search/hot/word"
		data, err := web.RequestDataWith(web.NewDefaultClient(), url, "GET", "", web.RandUA(), nil)
		if err != nil {
			msg := message.Text("ERROR:", err)
			ctx.SendChain(msg)
			return
		}
		json := gjson.Get(binary.BytesToString(data), "data").Array()
		for idx, hot := range json {
			if hot.Get("word").String() == "" {
				continue
			}
			rsp = rsp + strconv.Itoa(idx+1) + ":" + hot.Get("word").String() + "\n"
		}
		ctx.SendChain(message.Text(rsp))
	})

	engine.OnFullMatch("今日早报").Limit(ctxext.LimitByGroup).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getRealLink, err := bilibili.GetRealURL("https://xvfr.com/60s.php")
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR:", err))
			return
		}
		ctx.Send(message.Image(getRealLink))
	})
}
