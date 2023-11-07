// Package score 简单的积分系统
package score

import (
	"encoding/json"
	"math"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"

	coins "github.com/FloatTech/ZeroBot-Plugin/compounds/coins"

	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/zbputils/control"
	"github.com/tidwall/gjson"

	fcext "github.com/FloatTech/floatbox/ctxext"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type partygame struct {
	Name  string `json:"name"`
	Desc  string `json:"desc"`
	Coins string `json:"coins"`
}

var (
	pgs            = make(pg, 256)
	RobTimeManager = rate.NewManager[int64](time.Minute*70, 163)
	checkLimit     = rate.NewManager[int64](time.Minute*1, 5) // time setup
	catchLimit     = rate.NewManager[int64](time.Hour*1, 9)   // time setup
	processLimit   = rate.NewManager[int64](time.Hour*1, 5)   // time setup
	payLimit       = rate.NewManager[int64](time.Hour*1, 10)  // time setup
	wagerData      map[string]int
)

type pg = map[string]partygame

func init() {
	wagerData = make(map[string]int)
	wagerData["data"] = rand.Intn(2000)
	sdb := coins.Initialize("./data/score/score.db")
	loadFiles := fcext.DoOnceOnSuccess(func(ctx *zero.Ctx) bool {
		data, err := os.ReadFile(engine.DataFolder() + "loads.json")
		if err != nil {
			ctx.SendChain(message.Text("ERROR:", err))
			return false
		}
		err = json.Unmarshal(data, &pgs)
		if err != nil {
			panic(err)
		}
		return true
	})

	// 借鉴了其他bot的功能 编写而成
	engine.OnFullMatch("柠檬片总数", zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		uid := ctx.Event.UserID
		si := coins.GetSignInByUID(sdb, uid)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("您的柠檬片数量一共是: ", si.Coins))
	})
	engine.OnRegex(`^查询(\[CQ:at,qq=(\d+)\]\s?|(\d+))的柠檬片`, zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		TargetInt, _ := strconv.ParseInt(ctx.State["regex_matched"].([]string)[2], 10, 64)
		siTargetUser := coins.GetSignInByUID(sdb, TargetInt)
		getTargetName := ctx.CardOrNickName(TargetInt)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("这位 ( ", getTargetName, " ) 的柠檬片为", siTargetUser.Coins, "个"))
	})
	engine.OnFullMatch("抽奖", loadFiles, zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		if !checkLimit.Load(ctx.Event.UserID).Acquire() {
			ctx.SendChain(message.Text("太贪心了哦~过会试试吧"))
			return
		}
		getProtectStatus := CheckUserIsEnabledProtectMode(ctx.Event.UserID, sdb)
		if getProtectStatus {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经启动保护模式，不允许参与任何抽奖性质类互动"))
			return
		}
		var mutex sync.RWMutex // 添加读写锁以保证稳定性
		mutex.Lock()
		uid := ctx.Event.UserID
		si := coins.GetSignInByUID(sdb, uid) // 获取用户目前状况信息
		userCurrentCoins := si.Coins         // loading coins status
		if userCurrentCoins < 0 {
			_ = coins.InsertUserCoins(sdb, uid, 0)
			ctx.SendChain(message.Reply(uid), message.Text("本次参与的柠檬片不够哦~请多多打卡w"))
			return
		} // fix unexpected bug during the code error
		checkEnoughCoins := coins.CheckUserCoins(userCurrentCoins)
		if !checkEnoughCoins {
			ctx.SendChain(message.Reply(uid), message.Text("本次参与的柠檬片不够哦~请多多打卡w"))
			return
		}
		all := rand.Intn(39) // 一共39种可能性
		referpg := pgs[(strconv.Itoa(all))]
		getName := referpg.Name
		getCoinsStr := referpg.Coins
		getCoinsInt, _ := strconv.Atoi(getCoinsStr)
		getDesc := referpg.Desc
		addNewCoins := si.Coins + getCoinsInt - 60
		_ = coins.InsertUserCoins(sdb, uid, addNewCoins)
		msgid := ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(" 嗯哼~来玩抽奖了哦w 看看能抽到什么呢w"))
		time.Sleep(time.Second * 3)
		ctx.SendChain(message.Reply(msgid), message.Text("呼呼~让咱看看你抽到了什么东西ww\n"),
			message.Text("你抽到的是~ ", getName, "\n", "获得了柠檬片 ", getCoinsInt, "\n", getDesc, "\n目前的柠檬片总数为：", addNewCoins))
		mutex.Unlock()
	})
	// 一次最多骗 400 柠檬片,失败概率较大,失败会被反吞柠檬片
	engine.OnRegex(`^抢(\[CQ:at,qq=(\d+)\]\s?|(\d+))的柠檬片`, zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		if !catchLimit.Load(ctx.Event.UserID).Acquire() {
			ctx.SendChain(message.Text("太贪心了哦~一小时后再来试试吧"))
			return
		}
		getProtectStatus := CheckUserIsEnabledProtectMode(ctx.Event.UserID, sdb)
		if getProtectStatus {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经启动保护模式，不允许参与任何抽奖性质类互动"))
			return
		}
		TargetInt, _ := strconv.ParseInt(ctx.State["regex_matched"].([]string)[2], 10, 64)
		if TargetInt == ctx.Event.UserID {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("哈? 干嘛骗自己的?坏蛋哦"))
			return
		}
		getProtectTargetStatus := CheckUserIsEnabledProtectMode(TargetInt, sdb)
		if getProtectTargetStatus {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("对方已经启动保护模式，不允许此处理操作"))
			return
		}
		uid := ctx.Event.UserID
		siEventUser := coins.GetSignInByUID(sdb, uid)        // 获取主用户目前状况信息
		siTargetUser := coins.GetSignInByUID(sdb, TargetInt) // 获得被抢用户目前情况信息
		switch {
		case siEventUser.Coins < 400:
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("貌似没有足够的柠檬片去准备哦~请多多打卡w"))
			return
		case siTargetUser.Coins < 400:
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("太坏了~试图的对象貌似没有足够多的柠檬片~"))
			return
		}
		eventUserName := ctx.CardOrNickName(uid)
		eventTargetName := ctx.CardOrNickName(TargetInt)
		// token chance.
		// add more possibility to get the chance (0-200)
		getTicket := RobOrCatchLimitManager(ctx) // full is 1 , least 3. level 1,2,3,
		// however, the total is still 0-400.
		fullChanceToken := rand.Intn(10)
		var modifyCoins int
		if fullChanceToken > 7 { // use it to reduce the chance to lower coins.
			modifyCoins = rand.Intn(200) + 200
		} else {
			modifyCoins = rand.Intn(200)
		}
		getRandomNum := rand.Intn(10)
		PossibilityNum := 6 / getTicket
		setIsTrue := getRandomNum/PossibilityNum != 0
		var remindTicket string
		if getTicket == 3 {
			remindTicket = "目前已经达到疲倦状态，成功率下调到15%，或许考虑一下不要做一个坏人呢～ ^^ "
		}
		if setIsTrue {
			_ = coins.InsertUserCoins(sdb, siEventUser.UID, siEventUser.Coins-modifyCoins)
			_ = coins.InsertUserCoins(sdb, siTargetUser.UID, siTargetUser.Coins+modifyCoins)
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("试着去拿走 ", eventTargetName, " 的柠檬片时,被发现了.\n所以 ", eventUserName, " 失去了 ", modifyCoins, " 个柠檬片\n\n同时 ", eventTargetName, " 得到了 ", modifyCoins, " 个柠檬片\n", remindTicket))
			return
		}
		_ = coins.InsertUserCoins(sdb, siEventUser.UID, siEventUser.Coins+modifyCoins)
		_ = coins.InsertUserCoins(sdb, siTargetUser.UID, siTargetUser.Coins-modifyCoins)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("试着去拿走 ", eventTargetName, " 的柠檬片时,成功了.\n所以 ", eventUserName, " 得到了 ", modifyCoins, " 个柠檬片\n\n同时 ", eventTargetName, " 失去了 ", modifyCoins, " 个柠檬片\n", remindTicket))
	})
	engine.OnRegex(`^骗\s?\[CQ:at,qq=(\d+)\]\s(\d+)个柠檬片$`, zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		if !catchLimit.Load(ctx.Event.UserID).Acquire() {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("太贪心了哦~一小时后再来试试吧"))
			return
		}
		getProtectStatus := CheckUserIsEnabledProtectMode(ctx.Event.UserID, sdb)
		if getProtectStatus {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经启动保护模式，不允许参与任何抽奖性质类互动"))
			return
		}
		TargetInt, _ := strconv.ParseInt(ctx.State["regex_matched"].([]string)[1], 10, 64)
		getProtectTargetStatus := CheckUserIsEnabledProtectMode(TargetInt, sdb)
		if getProtectTargetStatus {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("对方已经启动保护模式，不允许此处理操作"))
			return
		}
		modifyCoins, _ := strconv.Atoi(ctx.State["regex_matched"].([]string)[2])
		if TargetInt == ctx.Event.UserID {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("哈? 干嘛骗自己的?坏蛋哦"))
			return
		}
		switch {
		case modifyCoins <= 100:
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("貌似你是想倒贴别人来着嘛?可以试试多骗一点哦，既然都骗了那就多点吧x"))
			return
		case modifyCoins > 2000:
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("不要太贪心了啦！太坏了 "))
			return
		}
		uid := ctx.Event.UserID
		siEventUser := coins.GetSignInByUID(sdb, uid)        // 获取主用户目前状况信息
		siTargetUser := coins.GetSignInByUID(sdb, TargetInt) // 获得被抢用户目前情况信息
		switch {
		case siTargetUser.Coins < modifyCoins:
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("太坏了~试图的对象貌似没有足够多的柠檬片~"))
			return
		case siEventUser.Coins < modifyCoins:
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("貌似你需要有那么多数量的柠檬片哦w"))
			return
		}
		eventUserName := ctx.CardOrNickName(uid)
		eventTargetName := ctx.CardOrNickName(TargetInt)
		// get random numbers.
		getTargetChanceToDealRaw := math.Round(float64(modifyCoins / 20)) // the total is 0-100，however I don't allow getting chance 0. lmao. max is 100 if modify is 2000
		getTicket := RobOrCatchLimitManager(ctx)
		var remindTicket string
		if getTicket == 3 {
			remindTicket = "目前已经达到疲倦状态，成功率下调本身概率的15%，或许考虑一下不要做一个坏人呢～ ^^ "
		}
		getTargetChanceToDealPossibilityKey := rand.Intn(102 / getTicket)
		if getTargetChanceToDealPossibilityKey < int(getTargetChanceToDealRaw) { // failed
			doubledModifyNum := modifyCoins * 2
			if doubledModifyNum > siEventUser.Coins {
				doubledModifyNum = siEventUser.Coins
				_ = coins.InsertUserCoins(sdb, siEventUser.UID, siEventUser.Coins-doubledModifyNum)
				_ = coins.InsertUserCoins(sdb, siTargetUser.UID, siTargetUser.Coins+doubledModifyNum)
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("试着去骗走 ", eventTargetName, " 的柠檬片时,被 ", eventTargetName, " 发现了.\n本该失去 ", modifyCoins*2, "\n但因为 ", eventUserName, " 的柠檬片过少，所以 ", eventUserName, " 失去了 ", doubledModifyNum, " 个柠檬片\n\n同时 ", eventTargetName, " 得到了 ", doubledModifyNum, " 个柠檬片\n", remindTicket))
				return
			}
			_ = coins.InsertUserCoins(sdb, siEventUser.UID, siEventUser.Coins-doubledModifyNum)
			_ = coins.InsertUserCoins(sdb, siTargetUser.UID, siTargetUser.Coins+doubledModifyNum)
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("试着去骗走 ", eventTargetName, " 的柠檬片时,被 ", eventTargetName, " 发现了.\n所以 ", eventUserName, " 失去了 ", doubledModifyNum, " 个柠檬片\n\n同时 ", eventTargetName, " 得到了 ", doubledModifyNum, " 个柠檬片\n", remindTicket))
			return
		}
		_ = coins.InsertUserCoins(sdb, siEventUser.UID, siEventUser.Coins+modifyCoins)
		_ = coins.InsertUserCoins(sdb, siTargetUser.UID, siTargetUser.Coins-modifyCoins)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("试着去拿走 ", eventTargetName, " 的柠檬片时,成功了.\n所以 ", eventUserName, " 得到了 ", modifyCoins, " 个柠檬片\n\n同时 ", eventTargetName, " 失去了 ", modifyCoins, " 个柠檬片\n", remindTicket))
	})
	engine.OnRegex(`^给\s?\[CQ:at,qq=(\d+)\]\s转(\d+)个柠檬片$`, zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		if !processLimit.Load(ctx.Event.UserID).Acquire() {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请等一会再转账哦w"))
			return
		}
		getProtectStatus := CheckUserIsEnabledProtectMode(ctx.Event.UserID, sdb)
		if getProtectStatus {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经启动保护模式，不允许参与任何抽奖性质类互动"))
			return
		}
		TargetInt, _ := strconv.ParseInt(ctx.State["regex_matched"].([]string)[1], 10, 64)
		getProtectTargetStatus := CheckUserIsEnabledProtectMode(TargetInt, sdb)
		if getProtectTargetStatus {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("对方已经启动保护模式，不允许此处理操作"))
			return
		}
		modifyCoins, _ := strconv.Atoi(ctx.State["regex_matched"].([]string)[2])
		if modifyCoins < 1 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("然而你不能转账低于0个柠檬片哦w～ 敲"))
			return
		}
		if TargetInt == ctx.Event.UserID {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("不可以给自己转账哦w（敲）"))
			return
		}
		uid := ctx.Event.UserID
		siEventUser := coins.GetSignInByUID(sdb, uid)        // 获取主用户目前状况信息
		siTargetUser := coins.GetSignInByUID(sdb, TargetInt) // 获得被转账用户目前情况信息
		if modifyCoins > siEventUser.Coins {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("貌似你的柠檬片数量不够哦~"))
			return
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("转账成功了哦~\n", ctx.CardOrNickName(siEventUser.UID), " 变化为 ", siEventUser.Coins, " - ", modifyCoins, "= ", siEventUser.Coins-modifyCoins, "\n", ctx.CardOrNickName(siTargetUser.UID), " 变化为: ", siTargetUser.Coins, " + ", modifyCoins, "= ", siTargetUser.Coins+modifyCoins))
		_ = coins.InsertUserCoins(sdb, siEventUser.UID, siEventUser.Coins-modifyCoins)
		_ = coins.InsertUserCoins(sdb, siTargetUser.UID, siTargetUser.Coins+modifyCoins)
	})
	engine.OnRegex(`^HandleCoins\s?\[CQ:at,qq=(\d+)\]\s(\d+)$`, zero.SuperUserPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		TargetInt, _ := strconv.ParseInt(ctx.State["regex_matched"].([]string)[1], 10, 64)
		modifyCoins, _ := strconv.Atoi(ctx.State["regex_matched"].([]string)[2])
		siTargetUser := coins.GetSignInByUID(sdb, TargetInt) // get user info
		unModifyCoins := siTargetUser.Coins
		_ = coins.InsertUserCoins(sdb, TargetInt, unModifyCoins+modifyCoins)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Handle Coins Successfully.\n"))
	})
	engine.OnRegex(`^(丢弃|扔掉)(\d+)个柠檬片$`).SetBlock(true).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		modifyCoins, _ := strconv.Atoi(ctx.State["regex_matched"].([]string)[2])
		handleUser := coins.GetSignInByUID(sdb, ctx.Event.UserID)
		currentUserCoins := handleUser.Coins
		if currentUserCoins-modifyCoins < 0 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("貌似你的柠檬片不够处理呢("))
			return
		}
		hadModifyCoins := currentUserCoins - modifyCoins
		_ = coins.InsertUserCoins(sdb, handleUser.UID, hadModifyCoins)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经帮你扔掉了哦"))
	})
	engine.OnFullMatch("兑换涩图", zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		if !payLimit.Load(ctx.Event.UserID).Acquire() {
			// 你 群 现 状
			ctx.SendChain(message.Text("坏欸！为什么一个群有这么多人看涩图啊（晕"))
			return
		}
		modified, _ := control.Lookup("nsfw")
		status := modified.IsEnabledIn(ctx.Event.GroupID)
		if status {
			var mutex sync.RWMutex // 添加读写锁以保证稳定性
			mutex.Lock()
			uid := ctx.Event.UserID
			si := coins.GetSignInByUID(sdb, uid) // 获取用户目前状况信息
			userCurrentCoins := si.Coins         // loading coins status
			if userCurrentCoins < 400 {
				ctx.SendChain(message.Reply(uid), message.Text("本次参与的柠檬片不够哦~请多多打卡w，一次兑换最少需要400"))
				return
			}
			img, err := web.GetData("https://api.lolicon.app/setu/v2?r18=1&num=1")
			if err != nil {
				ctx.SendChain(message.Text("ERROR:", err))
				return
			}
			picURL := gjson.Get(string(img), "data.0.urls.original").String()
			messageID := ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(picURL)).ID()
			if messageID != 0 { // 保证成功后才扣除
				_ = coins.InsertUserCoins(sdb, si.UID, userCurrentCoins-400)
			}
		} else {
			ctx.SendChain(message.Text("本群并没有开启nsfw哦，不允许使用此功能哦x"))
			return
		}
	})
	// I thought I just write a piece of shit. 💩
	engine.OnRegex(`^[!！]coin\swager\s?(\d*)`).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		// 得到本身奖池大小，如果没有或者被get的情况下获胜
		// this method should deal when we have less starter.
		rawNumber := ctx.State["regex_matched"].([]string)[1]
		if rawNumber == "" {
			rawNumber = "50"
		}
		getProtectStatus := CheckUserIsEnabledProtectMode(ctx.Event.UserID, sdb)
		if getProtectStatus {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经启动保护模式，不允许参与任何抽奖性质类互动"))
			return
		}
		modifyCoins, _ := strconv.Atoi(rawNumber)
		if modifyCoins > 1000 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("一次性最大投入为1k"))
			return
		}
		handleUser := coins.GetSignInByUID(sdb, ctx.Event.UserID)
		currentUserCoins := handleUser.Coins
		if currentUserCoins-modifyCoins < 0 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("貌似你的柠檬片不够处理呢("))
			return
		}
		// first of all , check the user status
		handlerWagerUser := coins.GetWagerUserStatus(sdb, ctx.Event.UserID)
		if handlerWagerUser.UserExistedStoppedTime > time.Now().Add(-time.Hour*12).Unix() {
			// then not pass | in the freeze time.
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("目前在冷却时间，距离下个可用时间为: ", time.Unix(handlerWagerUser.UserExistedStoppedTime, 0).Add(time.Hour*12).Format(time.DateTime)))
			return
		}
		// passed,delete this one and continue || before max is 3500.
		checkUserWagerCoins := handlerWagerUser.InputCountNumber
		if int64(modifyCoins)+checkUserWagerCoins > 3500 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("达到冷却最大值，您目前可投入："+strconv.Itoa(int(3500-checkUserWagerCoins))))
			return
		}
		// get wager
		getWager := coins.GetWagerStatus(sdb)
		if getWager.Expected == 0 {
			// it shows that no condition happened.
			// if not maxzine
			// in the wager mode. || start to load
			getGenOne := fcext.RandSenderPerDayN(time.Now().Unix(), 16500)
			getRandNumber := getGenOne + fcext.RandSenderPerDayN(time.Now().Unix()+ctx.Event.UserID, 5000) + 3000
			_ = coins.WagerCoinsInsert(sdb, modifyCoins+wagerData["data"], 0, getRandNumber)
			if int64(modifyCoins)+checkUserWagerCoins == 3500 {
				_ = coins.UpdateWagerUserStatus(sdb, ctx.Event.UserID, time.Now().Unix(), 0)
			} else {
				_ = coins.UpdateWagerUserStatus(sdb, ctx.Event.UserID, 0, int64(modifyCoins)+checkUserWagerCoins)
			}
			if getRandNumber <= modifyCoins {
				// winner, he | she is so lucky.^^
				// Lucy will cost 10 percent Coins.
				willRunCoins := math.Round(float64(modifyCoins+getWager.Wagercount) * 0.9)
				_ = coins.InsertUserCoins(sdb, ctx.Event.UserID, handleUser.Coins+int(willRunCoins)-modifyCoins)
				_ = coins.WagerCoinsInsert(sdb, 0, int(ctx.Event.UserID), 0)
				wagerData["data"] = int(math.Round(float64(modifyCoins+getWager.Wagercount)*0.1)) - 200
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("w！恭喜哦，奖池中奖了ww，一共获得 ", willRunCoins, " 个柠檬片，当前有 ", handleUser.Coins+int(willRunCoins)-modifyCoins, " 个柠檬片 (获胜者得到奖池 x0.9的柠檬片总数)"))
				return
			}
			// not winner
			_ = coins.InsertUserCoins(sdb, handleUser.UID, handleUser.Coins-modifyCoins)
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("没有中奖哦~，当前奖池为："+strconv.Itoa(modifyCoins)))
			return
		}
		// not init,start to add.
		getExpected := getWager.Expected
		if int64(modifyCoins)+checkUserWagerCoins == 3500 {
			_ = coins.UpdateWagerUserStatus(sdb, ctx.Event.UserID, time.Now().Unix(), 0)
		} else {
			_ = coins.UpdateWagerUserStatus(sdb, ctx.Event.UserID, 0, int64(modifyCoins)+checkUserWagerCoins)
		}
		if getWager.Wagercount+modifyCoins >= getExpected {
			// you are winner!
			willRunCoins := math.Round(float64(modifyCoins+getWager.Wagercount) * 0.9)
			_ = coins.InsertUserCoins(sdb, ctx.Event.UserID, handleUser.Coins+int(willRunCoins)-modifyCoins)
			_ = coins.WagerCoinsInsert(sdb, 0, int(ctx.Event.UserID), 0)
			wagerData["data"] = int(math.Round(float64(modifyCoins+getWager.Wagercount)*0.1)) - 200
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("w！恭喜哦，奖池中奖了ww，一共获得 ", willRunCoins, " 个柠檬片，当前有 ", handleUser.Coins+int(willRunCoins)-modifyCoins, " 个柠檬片 (获胜者得到奖池 x0.9的柠檬片总数)"))
			return
		} else {
			_ = coins.WagerCoinsInsert(sdb, getWager.Wagercount+modifyCoins, 0, getExpected)
			_ = coins.InsertUserCoins(sdb, ctx.Event.UserID, handleUser.Coins-modifyCoins)
			if rand.Intn(10) == 8 {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("呐～，不会还有大哥哥到现在 "+strconv.Itoa(getWager.Wagercount+modifyCoins)+" 个柠檬片了都没中奖吧？杂鱼～❤，杂鱼～❤"))
			} else {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("没有中奖哦~，当前奖池为: ", getWager.Wagercount+modifyCoins))
			}
		}
	})
	engine.OnRegex(`^(禁用|启用)柠檬片互动`).SetBlock(true).Limit(ctxext.LimitByUser).Handle(func(ctx *zero.Ctx) {
		getCode := ctx.State["regex_matched"].([]string)[1]
		if getCode == "禁用" {
			getStatus := CheckUserIsEnabledProtectMode(ctx.Event.UserID, sdb)
			if getStatus {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("你已经关闭了~"))
				return
			}
			// start to handle
			suser := coins.GetProtectModeStatus(sdb, ctx.Event.UserID)
			boolStatus := suser.Time+60*60*24 < time.Now().Unix()
			if !boolStatus {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("仅允许24小时修改一次"))
				return
			} // not the time
			// handle it.
			_ = coins.ChangeProtectStatus(sdb, ctx.Event.UserID, 1)
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("修改完成~"))
			return
		}
		getStatus := CheckUserIsEnabledProtectMode(ctx.Event.UserID, sdb)
		if !getStatus {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("你已经启用了~"))
			return
		}
		// start to handle
		suser := coins.GetProtectModeStatus(sdb, ctx.Event.UserID)
		boolStatus := suser.Time+60*60*24 < time.Now().Unix()
		if !boolStatus {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("仅允许24小时修改一次"))
			return
		} // not the time
		// handle it.
		_ = coins.ChangeProtectStatus(sdb, ctx.Event.UserID, 0)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("修改完成~"))
	})
}

func RobOrCatchLimitManager(ctx *zero.Ctx) (ticket int) {
	// use limitManager to reduce the chance of true.
	// 33 * 4 + 6 * 5 + null key (4 time tired.)
	/*
		first time to get full chance to win.
		second time reduce it to 50 % chance to win
		last time is null , you are tired and reduce it to 33% chance to win.
	*/
	switch {
	case RobTimeManager.Load(ctx.Event.UserID).AcquireN(33):
		return 1
	case RobTimeManager.Load(ctx.Event.UserID).AcquireN(6):
		return 2
	case RobTimeManager.Load(ctx.Event.UserID).Acquire():
		return 3
	default:
		return 3
	}
}

// CheckUserIsEnabledProtectMode 1 is enabled protect mode.
func CheckUserIsEnabledProtectMode(uid int64, sdb *coins.Scoredb) bool {
	s := coins.GetProtectModeStatus(sdb, uid)
	getCode := s.Status
	if getCode == 0 {
		return false
	} else {
		return true
	}
}
