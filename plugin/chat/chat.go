// Package chat basicchat
package chat

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/floatbox/process"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/FloatTech/ZeroBot-Plugin/compounds/name"
)

var (
	poke = rate.NewManager[int64](time.Minute*10, 8) // 戳一戳

	img    = "file:///root/Lucy_Project/memes/"
	engine = control.Register("chat", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "chat\n- [BOT名字]\n- [戳一戳BOT]\n- 空调开\n- 空调关\n- 群温度\n- 设置温度[正整数]",
	})
)

func init() { // 插件主体
	engine.OnRegex(`叫我.*?(.*)`, zero.OnlyToMe).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		texts := ctx.State["regex_matched"].([]string)[1]
		if name.StringInArray(texts, []string{"Lucy", "笨蛋", "老公", "猪", "夹子", "主人"}) {
			ctx.Send(message.Text("这些名字可不好哦(敲)"))
			return
		}

		if texts == "" {
			ctx.Send(message.Text("好哦~ 那~咱该叫你什么呢ww"))
			nextstep := ctx.FutureEvent("message", ctx.CheckSession())
			recv, cancel := nextstep.Repeat()
			for i := range recv {
				texts := i.MessageString()
				if texts != "" {
					cancel()
				}
			}
		}
		if strings.Contains(texts, "[CQ") {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("嗯哼 咱知道你在做什么哦"))
			return
		}
		userID := strconv.FormatInt(ctx.Event.UserID, 10)
		err := name.StoreUserNickname(userID, texts)
		if err != nil {
			ctx.Send(message.Text("发生了一些不可预料的问题 请稍后再试, ERR: ", err))
			return
		}
		ctx.Send(message.Text("好哦~ ", texts, " ちゃん~~~"))
	})

	// 被喊名字
	engine.OnFullMatch("", zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			var nickname = zero.BotConfig.NickName[0]
			time.Sleep(time.Second * 1)
			ctx.SendChain(message.Text(
				[]string{
					"这里是" + nickname + "(っ●ω●)っ",
					nickname + "不在呢~",
					"哼！" + nickname + "不想理你~",
				}[rand.Intn(3)],
			))
			process.SleepAbout1sTo2s()
			ctx.Send(message.Poke(ctx.Event.UserID))
		})

	// 戳一戳
	engine.On("notice/notify/poke", zero.OnlyToMe).SetBlock(false).
		Handle(func(ctx *zero.Ctx) {
			var nickname = zero.BotConfig.NickName[0]
			switch {
			case poke.Load(ctx.Event.GroupID).AcquireN(3):
				// 10分钟共8块命令牌 一次消耗3块命令牌
				time.Sleep(time.Second * 1)
				ctx.SendChain(message.Text([]string{"请不要戳" + nickname + " >_<", "再戳也不会理你的哦！", "别以为人家会搭理哦！",
					"呜…别戳了…", "别戳了！", "喵~", "有笨蛋在戳我，我不说是谁", "达咩呦，达咩达咩", "哼!不许戳啦 大笨蛋", "别戳啦！", "有笨蛋~让咱看看是谁"}[rand.Intn(11)]), message.Image([]string{img + "2941750127783.jpg", img + "C(185HMG2G0FY`3~2_[_H)W.gif", img + "file_3491851.jpg", img + "file_3492326.jpg", img + "file_3492330.jpg", img + "load.jpg"}[rand.Intn(6)]))
			case poke.Load(ctx.Event.GroupID).Acquire():
				// 5分钟共8块命令牌 一次消耗1块命令牌
				time.Sleep(time.Second * 1)
				happyFew := fmt.Sprintf("（好感 - %d）", rand.Intn(100)+1)
				ctx.SendChain(message.Text("喂(#`O′) 戳", nickname, "干嘛！", happyFew))
				process.SleepAbout1sTo2s()
				ctx.Send(message.Poke(ctx.Event.UserID))
			default:
				// 频繁触发，不回复
			}
		})

	// 戳我
	engine.OnFullMatchGroup([]string{"戳我", "戳戳"}, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			trueOrNot := rand.Intn(100)
			if trueOrNot >= 50 {
				process.SleepAbout1sTo2s()
				ctx.SendChain(message.At(ctx.Event.UserID), message.Text("好哦w"))
				process.SleepAbout1sTo2s()
				ctx.Send(message.Poke(ctx.Event.UserID))
			} else {
				process.SleepAbout1sTo2s()
				ctx.Send(message.Text("哼！Lucy才不想戳"))
			}
		})
	engine.OnKeywordGroup([]string{"会什么", "用法"}, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randImage("file_3492331.jpg"), RandText(
				"可以试着发送Lucy.help呢(",
			))
		})
	engine.OnFullMatch("捏脸", zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(RandText("大笨蛋！不许捏٩(๑`^´๑)۶", "疼....不许这样！哼！"), randImage("26329371069850.jpg", "2941750127783.jpg", "2OTN7BQ_1`YOPRH89[K{W8N.jpg"))
		})
	engine.OnFullMatch("摸头", zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randImage("6126814446620.jpg", "$OUXKWYM4LYHXT6)9I1WR5W.jpg", "file_3492330.jpg", "file_3491851.jpg", "file_3492333.jpg"), RandText("咱超可爱的w"))
		})
	engine.OnFullMatch("敲我", zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randImage("797198491dc98e4f.jpg", "file_3492319.jpg", "file_3492325.jpg"))
			process.SleepAbout1sTo2s()
			ctx.Send(randImage("6170420371656.jpg"))
			process.SleepAbout1sTo2s()
			ctx.SetGroupBan(
				ctx.Event.GroupID,
				ctx.Event.UserID,
				1*60)
			process.SleepAbout1sTo2s()
			ctx.SetGroupBan(
				ctx.Event.GroupID,
				ctx.Event.UserID,
				0)
		})
	engine.OnFullMatch("摸摸", zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(RandText("啾啾~", "呼呼~", "摸摸~"), randImage("8256CAEDA0E96A12875487BF2073256E.gif", "load.jpg", "-33ee3a0711f11810.jpg"))
		})
	engine.OnFullMatchGroup([]string{"呼呼", "抱抱"}, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randImage("1ce0f012eded2538.gif", "61KWD{AMBW[B3_AGSWJ6~}6.jpg", "wwwss.jpg"))
		})
	engine.OnFullMatch("抱住", zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randImage("3006028784945.jpg", "6126814446620.jpg", "2948706690280.jpg"))
		})
	engine.OnFullMatch("举高高", zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			process.SleepAbout1sTo2s()
			ctx.SendChain(message.Text("哼！才不让举高高呢"), randImage("file_3492332.jpg"))
		})
	engine.OnKeywordGroup([]string{"MoYoez", "MoeMagicMango"}, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(RandText("这些都是主人的名字~", "这些是架子(划掉)使用的名字w"))
		})
	engine.OnFullMatchGroup([]string{"是笨蛋"}, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			ctx.SendChain(randImage("SNDNYSG004[GH[E%$PJ~VCT.jpg", "55D0B4A5E335FE55A924E71469F35AC7.png", "file_3492326.jpg"))
		})
}

// RandText 随机文本
func RandText(text ...string) message.MessageSegment {
	return message.Text(text[rand.Intn(len(text))])
}

func randImage(file ...string) message.MessageSegment {
	return message.Image(img + file[rand.Intn(len(file))])
}

func GetPokeToken(ctx *zero.Ctx) float64 {
	return poke.Load(ctx.Event.GroupID).Tokens()
}
