package wife

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	coins "github.com/FloatTech/ZeroBot-Plugin/compounds/coins"
	"github.com/FloatTech/floatbox/binary"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/FloatTech/zbputils/img/text"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/extension/single"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	MessageTickerLimiter = rate.NewManager[int64](time.Minute*1, 2)
	engine               = control.Register("wife", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "Hi NekoPachi!",
		PrivateDataFolder: "wife",
	}).ApplySingle(ReverseSingle)
	ReverseSingle = single.New(
		single.WithKeyFn(func(ctx *zero.Ctx) int64 {
			return ctx.Event.UserID
		}),
		single.WithPostFn[int64](func(ctx *zero.Ctx) {
			if !MessageTickerLimiter.Load(ctx.Event.UserID).Acquire() {
				return
			}
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("正在操作哦～"))
		}),
	)
)

/*
StatusID:

Type 1: Normal Mode, nothing happened.

Type 2: Cannot be the Target, Target became initative, so reverse.

(However the target and the initative should be in their position, DO NOT CHANGE. )

Type 3: Something is wrong, you are Target == initative Person. (Drop The Person Before.)

Type 4: Removed.
(When User get others person. || IF REMARRIED, CHANGE IT TO TYPE1.) || (Be check more Time to reduce to err.)

Type 5: NTR Mode
(Tips: NTR means changed their pairkey & TargetID || UserID, need to do some changes. ) ||
(Attempt to do once more every person.)

Type 6: No wife Mod?
Fake - Invisibile person here.
(Lucy Hides this and shows it in the next Time if a person uses NTR,
shows nothing, and Lucy will make it for joke. LMAO)

Type 7: NTRED BY SOMEONE.
*/

func init() {
	sdb := coins.Initialize("./data/score/score.db")
	dict := make(map[string][]string) // this dict is used to reply
	// dict path.
	dict["block"] = []string{"嗯哼？貌似没有找到哦w", "再试试哦w，或许有帮助w", "运气不太好哦，想一下办法呢x"}
	dict["success"] = []string{"Lucky For You~", "恭喜哦ww~ ", "这边来恭喜一下哦w～", "貌似很成功的一次尝试呢w~"}
	dict["failed"] = []string{"今天的运气有一点背哦~这一次没有成功呢x", "_(:з」∠)_下次还有机会 抱抱w", "没关系哦，虽然失败了但还有机会呢x"}
	dict["ntr"] = []string{"嗯哼～这位还是成功了呢x", "aaa 好怪 不过还是让你通过了 ^^ "}
	dict["lost_failed"] = []string{"为什么要分呢? 让咱捏捏w", "太坏了啦！不许！"}
	dict["lost_success"] = []string{"好呢w 就这样呢(", "已经成功了哦w"}
	dict["hide_mode"] = []string{"哼哼～ 哼唧", "喵喵喵？！"}

	// main Class
	engine.OnRegex(`^(娶|嫁)(\[CQ:at,qq=(\d+)\])|^(娶|嫁)(Lucy)|^(娶|嫁)(.*)$`, zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		getFullRegexIndex := ctx.MessageString()
		// get Stat's situation.
		var fiancee int64
		var choice string
		caseDefault := ctx.State["regex_matched"].([]string)[1]
		switch {
		case caseDefault != "":
			choice = ctx.State["regex_matched"].([]string)[1]
			getID, _ := strconv.ParseInt(ctx.State["regex_matched"].([]string)[3], 10, 64)
			fiancee = getID
		case ctx.State["regex_matched"].([]string)[4] != "":
			choice = ctx.State["regex_matched"].([]string)[4]
			fiancee = ctx.Event.SelfID
		default:
			choice = ctx.State["regex_matched"].([]string)[6]
			getid := SearchUserReferName(ctx.Event.UserID, ctx.State["regex_matched"].([]string)[7])
			if getid == 0 {
				return
			}
			fiancee = getid
		}
		//
		if caseDefault != "" {
			// do split to text
			getSlice := strings.Split(getFullRegexIndex, " ")
			getAdditonName := getSlice[1]
			if strings.Contains(getAdditonName, "[CQ") {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("嗯哼 咱知道你在做什么哦"))
				return
			}
			if getAdditonName != "" {
				// save it to data.
				if getAdditonName == "Lucy" || getAdditonName == "群友" {
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("不可以使用 Lucy 或者 群友 这个名字"))
				} else {
					_ = SetUserReferName(ctx.Event.UserID, getAdditonName, fiancee)
					ctx.SendChain(message.Text("已经存入用户名字为:" + getAdditonName))
				}
			}
		}
		if fiancee == 80000000 || ctx.Event.UserID == 80000000 || fiancee == 0 {
			ctx.SendChain(message.Reply(ctx.Event.UserID), message.Text("用户不合法"))
			return
		}
		uid := ctx.Event.UserID
		if !CheckDisabledListIsExistedInThisGroup(marryList, uid, ctx.Event.GroupID) {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("你已经禁用了被随机，所以不可以参与娶群友哦w"))
			return
		}
		// fast check
		if !CheckTheUserStatusAndDoRepeat(ctx) {
			return
		}
		if !CheckTheTargetUserStatusAndDoRepeat(ctx, fiancee) {
			return
		}
		// check the target status.
		getStatusIfBannned := CheckTheUserIsInBlackListOrGroupList(fiancee, uid, ctx.Event.GroupID)
		/*
			disabled_Target
			blacklist_Target
		*/
		if getStatusIfBannned {
			// blocked.
			GlobalCDModelCost(ctx)
			getReply := dict["block"][rand.Intn(len(dict["block"]))]
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(getReply))
			return
		}
		if GlobalCDModelCostLeastReply(ctx) == 0 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("今天的机会已经使用完了哦～12小时后再来试试吧"))
			return
		}
		// check others.
		if uid == fiancee {
			switch rand.Intn(5) {
			case 1:
				GlobalCDModelCost(ctx)
				ctx.Send(ReplyMeantMode("貌似Lucy故意添加了 --force 的命令，成功了(笑 ", ctx.Event.UserID, 1, ctx))
				generatePairKey := GenerateMD5(ctx.Event.UserID, ctx.Event.UserID, ctx.Event.GroupID)
				err := InsertUserGlobalMarryList(marryList, ctx.Event.GroupID, ctx.Event.UserID, ctx.Event.UserID, 3, generatePairKey)
				if err != nil {
					panic(err)
				}
			default:
				GlobalCDModelCost(ctx)
				ctx.SendChain(message.Text("笨蛋！娶你自己干什么a"))
			}
			return
		}
		// However Lucy is only available to be married. LOL.
		if fiancee == ctx.Event.SelfID {
			// not work yet, so just the next path.
			if rand.Intn(100) > 90 {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("笨蛋！不准娶~ ama"))
				GlobalCDModelCost(ctx)
				return
			} else {
				// do it.
				GlobalCDModelCost(ctx)
				getSuccessMsg := dict["success"][rand.Intn(len(dict["success"]))]
				// normal mode. nothing happened.
				ctx.Send(ReplyMeantMode(getSuccessMsg, fiancee, 1, ctx))
				generatePairKey := GenerateMD5(ctx.Event.UserID, fiancee, ctx.Event.GroupID)
				_ = InsertUserGlobalMarryList(marryList, ctx.Event.GroupID, ctx.Event.UserID, fiancee, 1, generatePairKey)
				return
			}
		}
		switch {
		case choice == "娶":
			ResuitTheReferUserAndMakeIt(ctx, dict, uid, fiancee)
		case choice == "嫁":
			ResuitTheReferUserAndMakeIt(ctx, dict, fiancee, uid)
		}
	})
	engine.OnFullMatch("娶群友", zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		/*
			Work:
			- Check the User Status, if the user is 1 or 0 || 10 ,then pause and do this handler.
			- Choose a person, do something acciednt. (if the person had, pause and give one more chance.)
			- Check the banned or Disabled Status (To Target,if had,then stoppped it and give no chance. Others has checked itself too. )
			- add this key.
			- add more feature.
		*/
		/*
			TODO: HIDE MODE TYPE 6
		*/
		uid := ctx.Event.UserID
		gid := ctx.Event.GroupID
		// Check if Disabled this group.
		if !CheckDisabledListIsExistedInThisGroup(marryList, uid, gid) {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("你已经禁用了被随机，所以不可以参与娶群友哦w"))
			return
		}
		// fast check (Check User Itself.)
		if !CheckTheUserStatusAndDoRepeat(ctx) {
			return
		}

		ChooseAPerson := GetUserListAndChooseOne(ctx)
		// ok , go next. || before that we should check this person is in the lucky list?

		// Luck Path. (Only available in marry action.)
		getLuckyChance, getLuckyPeople, getLuckyTime := CheckTheOrderListAndBackDetailed(ctx.Event.UserID, ctx.Event.GroupID)
		getCurrentTime := time.Now().Unix()
		getLuckyTimeToInt64, _ := strconv.ParseInt(getLuckyTime, 10, 64)

		if getLuckyChance > 10 && getLuckyTimeToInt64 < getCurrentTime {
			if getLuckyTimeToInt64+(24*60*60) < getCurrentTime {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("貌似时间过去的太久了 这一次算是撤销掉了哦x,不过返回消耗的柠檬片，并且机会不变"))
				_ = RemoveOrderToList(marryList, ctx.Event.UserID, gid)
				getUserID := coins.GetSignInByUID(sdb, ctx.Event.UserID)
				_ = coins.InsertUserCoins(sdb, ctx.Event.UserID, getUserID.Coins+1000)
				return
			}
			getTargetStatusCode, _ := CheckTheUserIsTargetOrUser(marryList, ctx, getLuckyPeople) // 判断这个target是否已经和别人在一起了，同时判断Type3
			if getTargetStatusCode == -1 {
				// do this?
				getExistedToken := GlobalCDModelCostLeastReply(ctx)
				if getExistedToken == 0 {
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("今天的机会已经使用完了哦～12小时后再来试试吧，不过这边可以透露一下～许愿池已经实现了哦w～不过给你保留这次机会w"))
					return
				}
				// check the target status.
				getStatusIfBannned := CheckTheUserIsInBlackListOrGroupList(getLuckyPeople, uid, gid)
				if getStatusIfBannned {
					// blocked.
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("看起来挺倒霉的～貌似对方在许愿的过程中加入了黑名单x,或者对方已经禁用了对于本群的功能，只能无情删掉了哦x,不过这一次机会不会被浪费掉，并且会返回相应的柠檬片～"))
					getUserID := coins.GetSignInByUID(sdb, ctx.Event.UserID)
					_ = coins.InsertUserCoins(sdb, ctx.Event.UserID, getUserID.Coins+1000)
					_ = RemoveOrderToList(marryList, uid, gid)
					return
				}
				// success .
				GlobalCDModelCost(ctx)
				getSuccessMsg := dict["success"][rand.Intn(len(dict["success"]))]
				// normal mode. nothing happened.
				ctx.Send(ReplyMeantMode("许愿池生效～\n"+getSuccessMsg, getLuckyPeople, 1, ctx))
				generatePairKey := GenerateMD5(ctx.Event.UserID, getLuckyPeople, ctx.Event.GroupID)
				_ = InsertUserGlobalMarryList(marryList, ctx.Event.GroupID, ctx.Event.UserID, getLuckyPeople, 1, generatePairKey)
				_ = RemoveOrderToList(marryList, ctx.Event.UserID, gid)
				return
			} else {
				// target not -1 (has others.)
				// didn't do it.
				GlobalCDModelCost(ctx)
				_ = RemoveOrderToList(marryList, ctx.Event.UserID, gid)
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("抱歉哦～虽然已经使用了愿望池，不过仍然没有成功呢awa～"))
				// handle this chance but no cares
				return
			}
		}

		// Luck Path end.

		if !CheckTheTargetUserStatusAndDoRepeat(ctx, ChooseAPerson) {
			return
		}
		// check the target status.
		getStatusIfBannned := CheckTheUserIsInBlackListOrGroupList(ChooseAPerson, uid, gid)
		/*
			disabled_Target
			blacklist_Target
		*/
		if getStatusIfBannned {
			// blocked.
			GlobalCDModelCost(ctx)
			getReply := dict["block"][rand.Intn(len(dict["block"]))]
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(getReply))
			return
		}
		// go next. do something colorful, pls cost something.

		getExistedToken := GlobalCDModelCostLeastReply(ctx)
		if getExistedToken == 0 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("今天的机会已经使用完了哦～12小时后再来试试吧"))
			return
		}
		// one chance to get himself | herself
		if ChooseAPerson == ctx.Event.UserID {
			// status code 3
			GlobalCDModelCost(ctx)
			// drop target pls.
			ctx.Send(ReplyMeantMode("嗯哼哼～抽到了自己，然而 Lucy 还是将双方写成一个人哦w （笑w ", ctx.Event.UserID, 1, ctx))
			generatePairKey := GenerateMD5(ctx.Event.UserID, ctx.Event.UserID, ctx.Event.GroupID)
			_ = InsertUserGlobalMarryList(marryList, ctx.Event.GroupID, ctx.Event.UserID, ctx.Event.UserID, 3, generatePairKey)

		}
		returnNumber := GetSomeRanDomChoiceProps(ctx)
		switch {
		case returnNumber == 1:
			GlobalCDModelCost(ctx)
			getSuccessMsg := dict["success"][rand.Intn(len(dict["success"]))]
			// normal mode. nothing happened.
			ctx.Send(ReplyMeantMode(getSuccessMsg, ChooseAPerson, 1, ctx))
			generatePairKey := GenerateMD5(ctx.Event.UserID, ChooseAPerson, ctx.Event.GroupID)
			_ = InsertUserGlobalMarryList(marryList, ctx.Event.GroupID, ctx.Event.UserID, ChooseAPerson, 1, generatePairKey)
		case returnNumber == 2:
			GlobalCDModelCost(ctx)
			ctx.Send(ReplyMeantMode("貌似很奇怪哦～因为某种奇怪的原因～1变成了0,0变成了1", ChooseAPerson, 0, ctx))
			generatePairKey := GenerateMD5(ChooseAPerson, ctx.Event.UserID, ctx.Event.GroupID)
			_ = InsertUserGlobalMarryList(marryList, ctx.Event.GroupID, ChooseAPerson, ctx.Event.UserID, 2, generatePairKey)
		// reverse Target Mode
		case returnNumber == 3:
			GlobalCDModelCost(ctx)
			// drop target pls.
			ctx.Send(ReplyMeantMode("嗯哼哼～发生了一些错误～本来应当抽到别人的变成了自己～所以", ctx.Event.UserID, 1, ctx))
			generatePairKey := GenerateMD5(ctx.Event.UserID, ctx.Event.UserID, ctx.Event.GroupID)
			_ = InsertUserGlobalMarryList(marryList, ctx.Event.GroupID, ctx.Event.UserID, ctx.Event.UserID, 3, generatePairKey)
		// you became your own target
		case returnNumber == 6:
			GlobalCDModelCost(ctx)
			// now no wife mode.
			getHideMsg := dict["hide_mode"][rand.Intn(len(dict["hide_mode"]))]
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(getHideMsg, "\n貌似没有任何反馈～"))
			generatePairKey := GenerateMD5(ctx.Event.UserID, ChooseAPerson, ctx.Event.GroupID)
			_ = InsertUserGlobalMarryList(marryList, ctx.Event.GroupID, ctx.Event.UserID, ctx.Event.UserID, 6, generatePairKey)
		}
	})
	engine.OnFullMatch("离婚", zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		getStatusCode, _ := CheckTheUserIsTargetOrUser(marryList, ctx, ctx.Event.UserID)
		if getStatusCode == -1 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("貌似？没有对象的样子x"))
			return
		}
		if LeaveCDModelCostLeastReply(ctx) == 0 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("今天的次数已经用完了哦～或许可以试一下别的方式？"))
			return
		}
		reverseCheckTheUserIsDisabled := CheckDisabledListIsExistedInThisGroup(marryList, ctx.Event.UserID, ctx.Event.GroupID)
		if !reverseCheckTheUserIsDisabled {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("你已经禁用了被随机，所以不可以参与娶群友哦w"))
			return
		}
		getlostSuccessedMsg := dict["lost_success"][rand.Intn(len(dict["lost_success"]))]
		getLostFailedMsg := dict["lost_failed"][rand.Intn(len(dict["lost_failed"]))]
		if rand.Intn(4) >= 2 {
			LeaveCDModelCost(ctx)
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(getLostFailedMsg))
		} else {
			LeaveCDModelCost(ctx)
			getPairKey := CheckThePairKey(marryList, ctx.Event.UserID, ctx.Event.GroupID)
			RemoveUserGlobalMarryList(marryList, getPairKey, ctx.Event.GroupID)
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(getlostSuccessedMsg))
		}
	})

	engine.OnRegex(`^骗(\[CQ:at,qq=(\d+)\]\s?|(\d+))`, zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		fid := ctx.State["regex_matched"].([]string)
		fiancee, _ := strconv.ParseInt(fid[2]+fid[3], 10, 64)
		if fiancee == 0 || ctx.Event.UserID == 0 {
			return
		}
		uid := ctx.Event.UserID
		reverseCheckTheUserIsDisabled := CheckDisabledListIsExistedInThisGroup(marryList, ctx.Event.UserID, ctx.Event.GroupID)
		if !reverseCheckTheUserIsDisabled {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("你已经禁用了被随机，所以不可以参与娶群友哦w"))
			return
		}
		if fiancee == uid {
			ctx.SendChain(message.Text("要骗别人哦~为什么要骗自己呢x"))
			return
		}
		// this case should other people existed.
		getStatusCode, _ := CheckTheUserIsTargetOrUser(marryList, ctx, ctx.Event.UserID)
		if getStatusCode != -1 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("貌似你已经有了哦～？难不成时要找 ^^ 别人嘛（恼"))
			return
		}
		getTargetStatusCode, _ := CheckTheUserIsTargetOrUser(marryList, ctx, fiancee)
		if getTargetStatusCode == -1 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("嗯哼～这位还是一个人哦w～可以不用这个的哦w"))
			return
		}
		// low possibility to get this chance.
		if LeaveCDModelCostLeastReply(ctx) <= 0 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("今日机会不够哦w，过段时间再来试试吧w"))
			return
		}
		LeaveCDModelCost(ctx)
		if rand.Intn(100) < 30 {
			// win this goal
			getNTRMsg := dict["ntr"][rand.Intn(len(dict["ntr"]))]
			ctx.Send(ReplyMeantMode(getNTRMsg, fiancee, 5, ctx))
			CustomRemoveUserGlobalMarryList(marryList, CheckThePairKey(marryList, fiancee, ctx.Event.GroupID), ctx.Event.GroupID, 7)
			pairKey := GenerateMD5(ctx.Event.UserID, fiancee, ctx.Event.GroupID)
			err := InsertUserGlobalMarryList(marryList, ctx.Event.GroupID, ctx.Event.UserID, fiancee, 5, pairKey)
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
				return
			}
		} else {
			getFailed := dict["failed"][rand.Intn(len(dict["failed"]))]
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(getFailed))
			return
		}
	})
	engine.OnFullMatch("群老婆列表", zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		getList, num := GetTheGroupList(ctx.Event.GroupID)
		var RawMsg string
		if num == 0 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("本群貌似还没有人结婚来着（"))
			return
		}
		for i := 0; i < num; i++ {
			getUserInt64, _ := strconv.ParseInt(getList[i][0], 10, 64)
			getTargetInt64, _ := strconv.ParseInt(getList[i][1], 10, 64)
			RawMsg += strconv.FormatInt(int64(i+1), 10) + ". " + ctx.CardOrNickName(getUserInt64) + "( " + getList[i][0] + " )" + "  -->  " + ctx.CardOrNickName(getTargetInt64) + "( " + getList[i][1] + " )" + "\n"
		}
		headerMsg := "群老婆列表～ For Group( " + strconv.FormatInt(ctx.Event.GroupID, 10) + " )" + " [ " + ctx.GetGroupInfo(ctx.Event.GroupID, false).Name + " ]\n\n"
		base64Font, err := text.RenderToBase64(headerMsg+RawMsg+"\n\n Tips: 此列表将会在 23：00 PM (GMT+8) 重置", text.BoldFontFile, 1920, 45)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
			return
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image("base64://"+binary.BytesToString(base64Font)))
	})
	engine.OnRegex(`^添加黑名单(\[CQ:at,qq=(\d+)\]\s?|(\d+))`).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		fid := ctx.State["regex_matched"].([]string)
		fiancee, _ := strconv.ParseInt(fid[2]+fid[3], 10, 64)
		_ = AddBlackList(marryList, ctx.Event.UserID, fiancee)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经加入了～"))
	})
	engine.OnRegex(`^移除黑名单(\[CQ:at,qq=(\d+)\]\s?|(\d+))`).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		fid := ctx.State["regex_matched"].([]string)
		fiancee, _ := strconv.ParseInt(fid[2]+fid[3], 10, 64)
		_ = DeleteBlackList(marryList, ctx.Event.UserID, fiancee)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经移除了～"))
	})
	engine.OnFullMatch("添加本群禁用群老婆", zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		gid := ctx.Event.GroupID
		_ = AddDisabledList(marryList, ctx.Event.UserID, gid)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经加入了～,在本群你将不会加入此游戏"))
	})
	engine.OnFullMatch("删除本群禁用群老婆", zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		gid := ctx.Event.GroupID
		_ = DeleteDisabledList(marryList, ctx.Event.UserID, gid)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经移除了～,在本群你将加入此游戏"))
	})
	engine.OnRegex(`^添加许愿(\[CQ:at,qq=(\d+)\]\s?|(\d+))`).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		fid := ctx.State["regex_matched"].([]string)
		si := coins.GetSignInByUID(sdb, ctx.Event.UserID)
		fiancee, _ := strconv.ParseInt(fid[2]+fid[3], 10, 64)
		reverseCheckTheUserIsDisabled := CheckDisabledListIsExistedInThisGroup(marryList, ctx.Event.UserID, ctx.Event.GroupID)
		if !reverseCheckTheUserIsDisabled {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("你已经禁用了被随机，所以不可以参与娶群友哦w"))
			return
		}
		if si.Coins < 1000 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("本次许愿的柠檬片不足哦～需要1000个柠檬片才可以哦w"))
			return
		}
		if !CheckTheBlackListIsExistedToThisPerson(marryList, fiancee, ctx.Event.UserID) || !CheckTheBlackListIsExistedToThisPerson(marryList, ctx.Event.UserID, fiancee) {
			ctx.SendChain(message.Reply(ctx.Event.Message), message.Text("已经被Ban掉了，愿望无法实现～"))
			return
		}
		_, getTargetID, _ := CheckTheOrderListAndBackDetailed(ctx.Event.UserID, ctx.Event.GroupID)
		if getTargetID != fiancee && getTargetID != 0 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("每次仅可以许愿一个人w 不允许第二个人"))
			return
		}
		if fiancee == ctx.Event.UserID {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("坏哦！为什么要许自己的x"))
			return
		}
		if getTargetID == fiancee {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经许过一次了哦～不需要第二次"))
			return
		}
		// Handler
		_ = coins.InsertUserCoins(sdb, ctx.Event.UserID, si.Coins-1000)
		timeStamp := time.Now().Unix() + (6 * 60 * 60)
		_ = AddOrderToList(marryList, ctx.Event.UserID, fiancee, strconv.FormatInt(timeStamp, 10), ctx.Event.GroupID)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经许愿成功了哦～w 给", ctx.CardOrNickName(fiancee), " 的许愿已经生效，将会在6小时后增加70%可能性实现w"))
	})
	engine.OnFullMatch("重置娶群友", zero.SuperUserPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ResetToInitalizeMode()
		ctx.SendChain(message.Text("done"))
	})
}
