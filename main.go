// Package main for Lucy(HiMoYo Version)
package main

import (
	"flag"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"github.com/FloatTech/ZeroBot-Plugin/kanban"           // 在最前打印 banner
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/bottle"  // 漂流瓶
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/manager" // 群管
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/nsfw"    // nsfw
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/tools"   // 工具

	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/atri"  // atri
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/mai"   // mai
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/quote" // quote

	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/action" // action For Lucy触发

	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/funwork" // 好玩的整合工具

	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/score" // 签到

	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/chat" // 回复

	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/wife" // wife plugin

	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/pgr" // pgr plugin

	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/whitelist" // whitelist plugin

	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/chan"   // chan plugin
	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/choose" // slash plugin

	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/slash" // choose

	_ "github.com/FloatTech/ZeroBot-Plugin/plugin/simai" // simai plugin should be the last it will cause command block.

	"github.com/FloatTech/floatbox/process"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type zbpcfg struct {
	Z zero.Config        `json:"zero"`
	W []*driver.WSClient `json:"ws"`
}

var config zbpcfg

func init() {
	// 解析命令行参数
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	// make it easy to write config rather than change the source code.
	sus := make([]int64, 0, 16)
	// 直接写死 AccessToken 时，请更改下面第二个参数
	token := flag.String("t", "", "Set AccessToken of WSClient.")
	// 直接写死 URL 时，请更改下面第二个参数
	url := flag.String("u", "ws://127.0.0.1:6700", "Set Url of WSClient.")
	// 默认昵称
	adana := flag.String("n", os.Getenv("name"), "Set default nickname.")
	prefix := flag.String("p", "/", "Set command prefix.")
	flag.Parse()

	for _, s := range flag.Args() {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			continue
		}
		sus = append(sus, i)
	}
	// 通过代码写死的方式添加主人账号
	sus = append(sus, 1292581422)
	config.W = []*driver.WSClient{driver.NewWebSocketClient(*url, *token)}
	config.Z = zero.Config{
		NickName:       append([]string{*adana}, "Lucy", "lucy", "Lucy酱"),
		CommandPrefix:  *prefix,
		SuperUsers:     sus,
		Driver:         []zero.Driver{config.W[0]},
		MaxProcessTime: time.Minute * 4,
		RingLen:        0,
	}
}

func main() {
	zero.OnFullMatchGroup([]string{".help", "/help"}).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(message.Text(kanban.Banner))
		})
	zero.RunAndBlock(&config.Z, process.GlobalInitMutex.Unlock)
}
