// Package nsfw for nsfw (
package nsfw

import (
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	engine = control.Register("nsfw", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: true,
		Help:             "Hi NekoPachi!\n",
	})
)

func init() {
	engine.OnFullMatch("Hello", zero.OnlyToMe).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("World!"))
	})
}
