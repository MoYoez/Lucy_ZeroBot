package choose

import (
	"math/rand"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	engine := control.Register("choose", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  true,
		Help:              "choose - 帮助做选择",
		PrivateDataFolder: "choose",
	})
	engine.OnRegex(`^(.*)还是(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getRegexF := ctx.State["regex_matched"].([]string)[1]
		getRegexS := ctx.State["regex_matched"].([]string)[2]
		if getRegexS == getRegexF {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("貌似你已经选好了x Lucy不用帮你做选择了awwww"))
			return
		}
		getRand := rand.Intn(2)
		var result string
		switch {
		case getRand == 0:
			result = getRegexF
		case getRand == 1:
			result = getRegexS
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("aww 看起来选择"+result+"会比较合适("))
	})
}
