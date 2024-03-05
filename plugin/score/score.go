// Package score 简单的积分系统
package score

import (
	"bytes"
	"encoding/base64"
	"github.com/FloatTech/floatbox/math"
	"image/jpeg"
	"math/rand"
	"strconv"
	"sync"
	"time"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"

	coins "github.com/FloatTech/ZeroBot-Plugin/compounds/coins"
	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/gg"
	_ "github.com/FloatTech/sqlite" // import sql
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	engine = control.Register("score", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "Hi NekoPachi!",
		PrivateDataFolder: "score",
	})
	mutex sync.Mutex // 添加读写锁以保证稳定性
)

func init() {
	cachePath := engine.DataFolder() + "scorecache/"
	sdb := coins.Initialize("./data/score/score.db")
	ctxext.SetDefaultLimiterManagerParam(time.Second*10, 2)

	engine.OnFullMatchGroup([]string{"签到", "打卡"}, zero.OnlyGroup).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			uid := ctx.Event.UserID
			// save time data by add 30mins (database save it, not to handle it when it gets ready.)
			// just handle data time when it on,make sure to decrease 30 mins when render the page(
			getNowUnixFormatElevenThirten := time.Now().Add(time.Minute * 30).Format("20060102")
			getStatus, _ := coins.GetUserIsSignInToday(sdb, uid)
			si := coins.GetSignInByUID(sdb, uid)
			drawedFile := cachePath + strconv.FormatInt(uid, 10) + getNowUnixFormatElevenThirten + "signin.png"

			if getStatus && si.Count != 0 {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("w~ 你今天已经签到过了哦w"))
				if file.IsExist(drawedFile) {
					ctx.SendChain(message.Image("file:///" + file.BOTPATH + "/" + drawedFile))
				}
				return
			}
			coinsGet := 300 + rand.Intn(200)
			coins.UpdateUserSignInValue(sdb, uid)
			_ = coins.InsertUserCoins(sdb, uid, si.Coins+coinsGet)
			_ = coins.InsertOrUpdateSignInCountByUID(sdb, uid, si.Count+1) // 柠檬片获取
			score := coins.GetScoreByUID(sdb, uid).Score
			score++
			_ = coins.InsertOrUpdateScoreByUID(sdb, uid, score)
			CurrentCountTable := coins.GetCurrentCount(sdb, getNowUnixFormatElevenThirten)
			handledTodayNum := CurrentCountTable.Counttime + 1
			_ = coins.UpdateUserTime(sdb, handledTodayNum, getNowUnixFormatElevenThirten)
			if time.Now().Hour() > 6 && time.Now().Hour() < 19 {
				// package for test draw.
				getTimeReplyMsg := coins.GetHourWord(time.Now()) // get time and msg
				currentTime := time.Now().Format("2006-01-02 15:04:05")
				// time day.
				dayTimeImg, _ := gg.LoadImage(engine.DataFolder() + "BetaScoreDay.png")
				dayGround := gg.NewContext(1920, 1080)
				dayGround.DrawImage(dayTimeImg, 0, 0)
				_ = dayGround.LoadFontFace(engine.DataFolder()+"dyh.ttf", 60)
				dayGround.SetRGB(0, 0, 0)
				// draw something with cautions Only (
				dayGround.DrawString(currentTime, 1270, 950)            // draw time
				dayGround.DrawString(getTimeReplyMsg, 50, 930)          // draw text.
				dayGround.DrawString(ctx.CardOrNickName(uid), 310, 110) // draw name :p why I should do this???
				_ = dayGround.LoadFontFace(engine.DataFolder()+"dyh.ttf", 60)
				dayGround.DrawStringWrapped(strconv.Itoa(handledTodayNum), 350, 255, 1, 1, 0, 1.3, gg.AlignCenter)   // draw first part
				dayGround.DrawStringWrapped(strconv.Itoa(si.Count+1), 1000, 255, 1, 1, 0, 1.3, gg.AlignCenter)       // draw second part
				dayGround.DrawStringWrapped(strconv.Itoa(coinsGet), 220, 370, 1, 1, 0, 1.3, gg.AlignCenter)          // draw third part
				dayGround.DrawStringWrapped(strconv.Itoa(si.Coins+coinsGet), 720, 370, 1, 1, 0, 1.3, gg.AlignCenter) // draw forth part
				// level array with rectangle work.
				rankNum := coins.GetLevel(score)
				RankGoal := rankNum + 1
				achieveNextGoal := coins.LevelArray[RankGoal]
				achievedGoal := coins.LevelArray[rankNum]
				currentNextGoalMeasure := achieveNextGoal - score            // measure rest of the num. like 20 - currentLink(TestRank 15)
				measureGoalsLens := math.Abs(achievedGoal - achieveNextGoal) // like 20 - 10
				currentResult := float64(measureGoalsLens-currentNextGoalMeasure) / float64(measureGoalsLens)
				// draw this part
				dayGround.SetRGB255(180, 255, 254)        // aqua color
				dayGround.DrawRectangle(70, 570, 600, 50) // draw rectangle part1
				dayGround.Fill()
				dayGround.SetRGB255(130, 255, 254)
				dayGround.DrawRectangle(70, 570, 600*currentResult, 50) // draw rectangle part2
				dayGround.Fill()
				dayGround.SetRGB255(0, 0, 0)
				dayGround.DrawString("Lv. "+strconv.Itoa(rankNum)+" 签到天数 + 1", 80, 490)
				_ = dayGround.LoadFontFace(engine.DataFolder()+"dyh.ttf", 40)
				dayGround.DrawString(strconv.Itoa(measureGoalsLens-currentNextGoalMeasure)+"/"+strconv.Itoa(measureGoalsLens), 710, 610)
				tureResult := dayGround.Image()
				var buf bytes.Buffer
				err := jpeg.Encode(&buf, tureResult, nil)
				if err != nil {
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
					return
				}
				base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
				ctx.SendChain(message.At(uid), message.Text("[HiMoYoBot]签到成功\n"), message.Image("base64://"+base64Str))
				_ = dayGround.SavePNG(drawedFile)
			} else {
				// nightVision
				// package for test draw.
				getTimeReplyMsg := coins.GetHourWord(time.Now()) // get time and msg
				currentTime := time.Now().Format("2006-01-02 15:04:05")
				nightTimeImg, _ := gg.LoadImage(engine.DataFolder() + "BetaScoreNight.png")
				nightGround := gg.NewContext(1886, 1060)
				nightGround.DrawImage(nightTimeImg, 0, 0)
				_ = nightGround.LoadFontFace(engine.DataFolder()+"dyh.ttf", 60)
				nightGround.SetRGB255(255, 255, 255)
				// draw something with cautions Only (
				nightGround.DrawString(currentTime, 1360, 910)            // draw time
				nightGround.DrawString(getTimeReplyMsg, 60, 930)          // draw text.
				nightGround.DrawString(ctx.CardOrNickName(uid), 350, 140) // draw name :p why I should do this???
				_ = nightGround.LoadFontFace(engine.DataFolder()+"dyh.ttf", 60)
				nightGround.DrawStringWrapped(strconv.Itoa(handledTodayNum), 345, 275, 1, 1, 0, 1.3, gg.AlignCenter)   // draw first part
				nightGround.DrawStringWrapped(strconv.Itoa(si.Count+1), 990, 275, 1, 1, 0, 1.3, gg.AlignCenter)        // draw second part
				nightGround.DrawStringWrapped(strconv.Itoa(coinsGet), 225, 360, 1, 1, 0, 1.3, gg.AlignCenter)          // draw third part
				nightGround.DrawStringWrapped(strconv.Itoa(si.Coins+coinsGet), 720, 360, 1, 1, 0, 1.3, gg.AlignCenter) // draw forth part
				// level array with rectangle work.
				rankNum := coins.GetLevel(score)
				RankGoal := rankNum + 1
				achieveNextGoal := coins.LevelArray[RankGoal]
				achievedGoal := coins.LevelArray[rankNum]
				currentNextGoalMeasure := achieveNextGoal - score            // measure rest of the num. like 20 - currentLink(TestRank 15)
				measureGoalsLens := math.Abs(achievedGoal - achieveNextGoal) // like 20 - 10
				currentResult := float64(measureGoalsLens-currentNextGoalMeasure) / float64(measureGoalsLens)
				// draw this part
				nightGround.SetRGB255(49, 86, 157)          // aqua color
				nightGround.DrawRectangle(70, 570, 600, 50) // draw rectangle part1
				nightGround.Fill()
				nightGround.SetRGB255(255, 255, 255)
				nightGround.DrawRectangle(70, 570, 600*currentResult, 50) // draw rectangle part2
				nightGround.Fill()
				nightGround.SetRGB255(255, 255, 255)
				nightGround.DrawString("Lv. "+strconv.Itoa(rankNum)+" 签到天数 + 1", 80, 490)
				_ = nightGround.LoadFontFace(engine.DataFolder()+"dyh.ttf", 40)
				nightGround.DrawString(strconv.Itoa(measureGoalsLens-currentNextGoalMeasure)+"/"+strconv.Itoa(measureGoalsLens), 710, 610)
				tureResult := nightGround.Image()
				var buf bytes.Buffer
				err := jpeg.Encode(&buf, tureResult, nil)
				if err != nil {
					ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
					return
				}
				base64Str := base64.StdEncoding.EncodeToString(buf.Bytes())
				ctx.SendChain(message.At(uid), message.Text("[HiMoYoBot]签到成功\n"), message.Image("base64://"+base64Str))
				_ = nightGround.SavePNG(drawedFile)

			}

		})
}
