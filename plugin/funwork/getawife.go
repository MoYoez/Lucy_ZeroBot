// Package funwork 简单本地老婆
package funwork

import (
	"os"
	"regexp"
	"time"

	fcext "github.com/FloatTech/floatbox/ctxext"
	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	cachePath := engine.DataFolder() + "wife/"
	engine.OnFullMatch("抽老婆").SetBlock(true).Limit(ctxext.LimitByGroup).
		Handle(func(ctx *zero.Ctx) {
			wifes, err := os.ReadDir(cachePath)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			name := ctx.CardOrNickName(ctx.Event.UserID)
			n := fcext.RandSenderPerDayN(ctx.Event.UserID, len(wifes))
			wn := wifes[n].Name()
			reg := regexp.MustCompile(`[^.]+`)
			list := reg.FindAllString(wn, -1)
			deleteme := ctx.SendChain(
				message.Text(name, "さんが二次元で結婚するであろうヒロインは、", "\n"),
				message.Image("file:///"+file.BOTPATH+"/"+cachePath+wn),
				message.Text("\n【", list[0], "】です！"))
			time.Sleep(time.Second * 15)
			ctx.DeleteMessage(deleteme)
		})
}
