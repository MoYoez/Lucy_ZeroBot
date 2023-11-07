// Package action for Lucy
package action

import (
	"math/rand"
	"time"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"

	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	limit   = rate.NewManager[int64](time.Minute*10, 15)
	LucyImg = "file:///root/Lucy_Project/memes/" // LucyImg for Lucy的meme表情包地址
)

func init() {
	engine := control.Register("action", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "Lucy容易被动触发语言\n",
	})
	engine.OnFullMatchGroup([]string{"喵", "喵喵", "喵喵喵"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			switch rand.Intn(2) {
			case 0:
				ctx.SendChain(RandText("喵喵~", "喵w~"))
			case 1:
				ctx.SendChain(randImage("6152277811454.jpg", "meow.jpg", "file_3491851.jpg", "file_3492320.jpg"))
			}
		})
	engine.OnFullMatch("咕咕").SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			if !limit.Load(ctx.Event.GroupID).Acquire() {
				return
			}
			ctx.SendChain(RandText("炖了~鸽子都要恰掉w", "咕咕咕", "不许咕咕咕"))
		})
}

// RandText for random text
func RandText(text ...string) message.MessageSegment {
	return message.Text(text[rand.Intn(len(text))])
}

func randImage(file ...string) message.MessageSegment {
	return message.Image(LucyImg + file[rand.Intn(len(file))])
}
