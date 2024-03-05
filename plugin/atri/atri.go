/*
Package atri 本文件基于 https://github.com/Kyomotoi/ATRI
为 Golang 移植版，语料、素材均来自上述项目
本项目遵守 AGPL v3 协议进行开源
*/
package atri

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/floatbox/process"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	"github.com/FloatTech/ZeroBot-Plugin/compounds/name"
)

const (
	servicename = "atri"
)

func init() { // 插件主体
	engine := control.Register(servicename, &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help: "- Lucy醒醒\n- Lucy睡吧\n- 萝卜子\n- 喜欢 | 爱你 | 爱 | suki | daisuki | すき | 好き | 贴贴 | 老婆 | 亲一个 | mua\n" +
			"- 早安 | 早哇 | 早上好 | ohayo | 哦哈哟 | お早う | 早好 | 早 | 早早早\n" +
			"- 中午好 | 午安 | 午好\n- 晚安 | oyasuminasai | おやすみなさい | 晚好 | 晚上好\n- 高性能 | 太棒了 | すごい | sugoi | 斯国一 | よかった\n" +
			"- 没事 | 没关系 | 大丈夫 | 还好 | 不要紧 | 没出大问题 | 没伤到哪\n- 好吗 | 是吗 | 行不行 | 能不能 | 可不可以\n- 啊这\n- 我好了\n- ？ | ? | ¿\n" +
			"- 离谱\n- 答应我",
	})

	engine.OnFullMatchGroup([]string{"早安", "早哇", "早上好", "ohayo", "哦哈哟", "お早う", "早好", "早", "早早早"}, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			now := time.Now().Hour()
			process.SleepAbout1sTo2s()
			switch {
			case now < 6: // 凌晨
				ctx.SendChain(message.Reply(ctx.Event.MessageID), RandWithReplaceName(ctx,
					"zzzz......",
					"zzzzzzzz......",
					"zzz...好涩哦..zzz....",
					"别...不要..zzz..那..zzz..",
					"嘻嘻..zzz..呐~..zzzz..",
					"...zzz....哧溜哧溜....",
				))
			case now >= 6 && now < 9:
				ctx.SendChain(message.Reply(ctx.Event.MessageID), RandWithReplaceName(ctx,
					"啊......早上好...(哈欠)",
					"唔......吧唧...早上...哈啊啊~~~\n早上好......",
					"早上好......",
					"早上好呜......呼啊啊~~~~",
					"吧唧吧唧......怎么了...已经早上了么...",
					"早上好！",
					"......看起来像是傍晚，其实已经早上了吗？",
					"早上好......欸~~~脸好近呢",
					"早安吖~新的一天要继续加油哦w",
					"早安~祝你的坏心情被Lucy带走~~~(*/ω＼*)",
				))
			case now >= 9 && now < 18:
				ctx.SendChain(message.Reply(ctx.Event.MessageID), RandWithReplaceName(ctx,
					"哼！这个点还早啥，昨晚干啥去了！？",
					"熬夜了对吧熬夜了对吧熬夜了对吧？？？！",
					"是不是熬夜是不是熬夜是不是熬夜？！",
					"Lucy酱提醒~不要熬夜哟w~ 否则就和架子一样被打死（）",
					"哼！ 这个点起床 一定是熬夜了！(｀･∪･´)",
				))
			case now >= 18 && now < 24:
				ctx.SendChain(message.Reply(ctx.Event.MessageID), RandWithReplaceName(ctx,
					"早个啥？哼唧！Lucy都准备洗洗睡了！",
					"不是...你看看几点了，哼！",
					"晚上好哇",
					"是时区问题嘛~ 还是上夜班呢~ 总之记得对自己好一点哦w",
				))
			}
		})

	engine.OnFullMatchGroup([]string{"中午好", "午安", "午好"}, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			now := time.Now().Hour()
			if now > 11 && now < 15 { // 中午
				process.SleepAbout1sTo2s()
				ctx.SendChain(message.Reply(ctx.Event.MessageID), RandWithReplaceName(ctx,
					"午安w",
					"午觉要好好睡哦，Lucy会陪伴在你身旁的w",
					"嗯哼哼~睡吧，就像平常一样安眠吧~o(≧▽≦)o",
					"睡你午觉去！哼唧！！",
				))
			}
		})
	engine.OnFullMatchGroup([]string{"晚安", "oyasuminasai", "おやすみなさい", "晚好", "晚上好"}, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			now := time.Now().Hour()
			process.SleepAbout1sTo2s()
			switch {
			case now < 6: // 凌晨
				ctx.SendChain(message.Reply(ctx.Event.MessageID), RandWithReplaceName(ctx,
					"zzzz......",
					"zzzzzzzz......",
					"...zzz....哧溜哧溜....",
				))
			case now >= 6 && now < 11:
				ctx.SendChain(message.Reply(ctx.Event.MessageID), RandWithReplaceName(ctx,
					"？啊这",
					"亲，这边建议赶快去睡觉呢~~~",
					"? 你知道现在几点了嘛 \n还不快去睡觉~",
				))
			case now >= 11 && now < 15:
				ctx.SendChain(message.Reply(ctx.Event.MessageID), RandWithReplaceName(ctx,
					"午安w",
					"午觉要好好睡哦，Lucy会陪伴在你身旁的w",
					"嗯哼哼~睡吧，就像平常一样安眠吧~o(≧▽≦)o",
					"睡你午觉去！哼唧！！",
				))
			case now >= 15 && now < 19:
				ctx.SendChain(message.Reply(ctx.Event.MessageID), RandWithReplaceName(ctx,
					"难不成？？晚上不想睡觉？？现在休息",
					"就......挺离谱的...现在睡觉",
					"现在还是白天哦，睡觉还太早了",
				))
			case now >= 19 && now < 24:
				ctx.SendChain(message.Reply(ctx.Event.MessageID), RandWithReplaceName(ctx,
					"嗯哼哼~睡吧，就像平常一样安眠吧~o(≧▽≦)o",
					"......(打瞌睡)",
					"呼...呼...已经睡着了哦~...呼......",
					"......Lucy...Lucy会在这守着你的，请务必好好睡着",
				))
			}
		})
	engine.OnKeywordGroup([]string{"高性能", "太棒了", "すごい", "sugoi", "斯国一", "よかった"}, atriSleep, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			process.SleepAbout1sTo2s()
			ctx.SendChain(RandWithReplaceName(ctx,
				"当然，Lucy是高性能的嘛~！",
				"小事一桩，Lucy是高性能的嘛",
				"怎么样？还是Lucy比较高性能吧？",
				"哼哼！Lucy果然是高性能的呢！",
				"因为Lucy是高性能的嘛！嗯哼！",
				"因为Lucy是高性能的呢！",
				"哎呀~，Lucy可真是太高性能了",
				"正是，因为Lucy是高性能的",
				"是的。Lucy是高性能的嘛♪",
				"毕竟Lucy可是高性能的！",
				"嘿嘿，Lucy的高性能发挥出来啦♪",
				"Lucy果然是很高性能的机器人吧！",
				"是吧！谁叫Lucy这么高性能呢！哼哼！",
				"交给Lucy吧，有高性能的Lucy陪着呢",
				"呣......Lucy的高性能，毫无遗憾地施展出来了......",
			))
		})

	engine.OnKeywordGroup([]string{"没事", "没关系", "大丈夫", "还好", "不要紧", "没出大问题", "没伤到哪"}, atriSleep, zero.OnlyToMe).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			process.SleepAbout1sTo2s()
			ctx.SendChain(RandWithReplaceName(ctx,
				"当然，Lucy是高性能的嘛~！",
				"没事没事，因为Lucy是高性能的嘛！嗯哼！",
				"没事的，因为Lucy是高性能的呢！",
				"正是，因为Lucy是高性能的",
				"是的。Lucy是高性能的嘛♪",
				"毕竟Lucy可是高性能的！",
				"那种程度的事不算什么的。\n别看Lucy这样，Lucy可是高性能的",
				"没问题的，Lucy可是高性能的",
			))
		})
}

// atriSleep 凌晨0点到6点，ATRI 在睡觉，不回应任何请求
func atriSleep(*zero.Ctx) bool {
	if now := time.Now().Hour(); now >= 0 && now < 6 {
		return false
	}
	return true
}

// RandWithReplaceName 随机返回一条带有替换的文本 用于回复 同时替换掉 “你”
func RandWithReplaceName(ctx *zero.Ctx, text ...string) message.MessageSegment {
	getNum := rand.Intn(len(text))
	IDStr := strconv.FormatInt(ctx.Event.UserID, 10)
	needToReplace := text[getNum]
	getName := name.LoadUserNickname(IDStr)
	output := strings.ReplaceAll(needToReplace, "你", getName)
	return message.Text(output)
}
