package pgr // Package pgr hosted by Phigros-Library
import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/tidwall/gjson"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/gg"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/disintegration/imaging"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

// too lazy,so this way is to use thrift host server (Working on HiMoYo Cloud.) (replace: now use PUA API)

// update: use PhigrosUnlimitedAPI + Phigros Library as Maintainer.

type QuerySongDetailsGenerator struct {
	Status  bool `json:"status"`
	Content struct {
		Songid string `json:"songid"`
		Info   struct {
			Songname    string `json:"songname"`
			Composer    string `json:"composer"`
			Illustrator string `json:"illustrator"`
			ChartDetail struct {
				EZ struct {
					Rating  float64 `json:"rating"`
					Charter string  `json:"charter"`
				} `json:"EZ"`
				HD struct {
					Rating  float64 `json:"rating"`
					Charter string  `json:"charter"`
				} `json:"HD"`
				In struct {
					Rating  float64 `json:"rating"`
					Charter string  `json:"charter"`
				} `json:"In"`
				At struct {
					Rating  float64 `json:"rating"`
					Charter string  `json:"charter"`
				} `json:"At"`
				LevelList []float64 `json:"level_list"`
			} `json:"chartDetail"`
		} `json:"info"`
	} `json:"content"`
}

var (
	engine = control.Register("phigros", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "Hi NekoPachi!\n",
		PrivateDataFolder: "phi",
	})
	router = "http://localhost:51602"
)

func init() {
	engine.OnRegex(`^[! ！/]pgr\sbind\s(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		hash := ctx.State["regex_matched"].([]string)[1]
		userInfo := GetUserInfoTimeFromDatabase(ctx.Event.UserID)
		if userInfo+(12*60*60) > time.Now().Unix() {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("12小时内仅允许绑定一次哦"))
			return
		}
		indexReply := DecHashToRaw(hash)
		// get session.
		if indexReply == "" {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请前往 https://phi.lemonkoi.one 获取绑定码进行绑定"))
			return
		}
		getQQID, getSessionID := RawJsonParse(indexReply)
		if getQQID != ctx.Event.UserID {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请求Hash中QQ号不一致，请使用自己的号重新申请"))
			return
		}
		if utf8.RuneCountInString(getSessionID) != 25 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Session 传入数值出现错误，请重新绑定"))
			return
		}
		_ = FormatUserDataBase(getQQID, getSessionID, time.Now().Unix()).BindUserDataBase()
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("绑定成功～"))
	})
	engine.OnRegex(`^[! ！/]pgr\sb19$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		data := GetUserInfoFromDatabase(ctx.Event.UserID)
		getDataSession := data.PhiSession
		if getDataSession == "" {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请前往 https://phi.lemonkoi.one 获取绑定码进行绑定 "))
			return
		}
		userData := GetUserInfoFromDatabase(ctx.Event.UserID)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("好哦~正在帮你请求，请稍等一下啦w~大约需要1-2分钟"))
		var dataWaiter sync.WaitGroup
		var AvatarWaiter sync.WaitGroup
		var getAvatarFormat *gg.Context
		var phidata []byte
		var setGlobalStat = true
		AvatarWaiter.Add(1)
		dataWaiter.Add(1)
		go func() {
			defer dataWaiter.Done()
			getFullLink := router + "/api/phi/bests?session=" + userData.PhiSession + "&overflow=2"
			phidata, _ = web.GetData(getFullLink)
			if phidata == nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("目前 Unoffical Phigros API 暂时无法工作 请过一段时候尝试"))
				setGlobalStat = false
				return
			}
			err := json.Unmarshal(phidata, &phigrosB19)
			if !phigrosB19.Status || err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("w? 貌似出现了一些问题x"))
				return
			}
		}()
		go func() {
			defer AvatarWaiter.Done()
			avatarByte, err := http.Get("https://q4.qlogo.cn/g?b=qq&nk=" + strconv.FormatInt(ctx.Event.UserID, 10) + "&s=640")
			if err != nil {
				ctx.SendChain(message.Text("Something wrong while rendering pic? avatar IO err."))
				return
			}
			// avatar
			avatarByteUni, _, _ := image.Decode(avatarByte.Body)
			showUserAvatar := imaging.Resize(avatarByteUni, 250, 250, imaging.Lanczos)
			getAvatarFormat = gg.NewContext(250, 250)
			getAvatarFormat.DrawRoundedRectangle(0, 0, 248, 248, 20)
			getAvatarFormat.Clip()
			getAvatarFormat.DrawImage(showUserAvatar, 0, 0)
			getAvatarFormat.Fill()
		}()
		getRawBackground, _ := gg.LoadImage(backgroundRender)
		getMainBgRender := gg.NewContextForImage(imaging.Resize(getRawBackground, 2750, 5500, imaging.Lanczos))
		_ = getMainBgRender.LoadFontFace(font, 30)
		// header background
		drawTriAngle(getMainBgRender, a, 0, 166, 1324, 410)
		getMainBgRender.SetRGBA255(0, 0, 0, 160)
		getMainBgRender.Fill()
		drawTriAngle(getMainBgRender, a, 1318, 192, 1600, 350)
		getMainBgRender.SetRGBA255(0, 0, 0, 160)
		getMainBgRender.Fill()
		drawTriAngle(getMainBgRender, a, 1320, 164, 6, 414)
		getMainBgRender.SetColor(color.White)
		getMainBgRender.Fill()
		// header background end.
		// load icon with other userinfo.
		getMainBgRender.SetColor(color.White)
		logo, _ := gg.LoadPNG(icon)
		getImageLogo := imaging.Resize(logo, 290, 290, imaging.Lanczos)
		getMainBgRender.DrawImage(getImageLogo, 50, 216)
		fontface, _ := gg.LoadFontFace(font, 90)
		getMainBgRender.SetFontFace(fontface)
		getMainBgRender.DrawString("Phigros", 422, 336)
		getMainBgRender.DrawString("RankingScore查询", 422, 462)
		// draw userinfo path
		renderHeaderText, _ := gg.LoadFontFace(font, 54)
		getMainBgRender.SetFontFace(renderHeaderText)
		dataWaiter.Wait()
		if !setGlobalStat {
			return
		}
		getMainBgRender.DrawString("Player: "+phigrosB19.Content.PlayerID, 1490, 300)
		getMainBgRender.DrawString("RankingScore: "+strconv.FormatFloat(phigrosB19.Content.RankingScore, 'f', 3, 64), 1490, 380)
		getMainBgRender.DrawString("ChanllengeMode: ", 1490, 460) // +56
		getColor, getLink := GetUserChallengeMode(phigrosB19.Content.ChallengeModeRank)
		if getColor != "" {
			getColorLink := ChanllengeMode + getColor + ".png"
			getColorImage, _ := gg.LoadImage(getColorLink)
			getMainBgRender.DrawImage(imaging.Resize(getColorImage, 238, 130, imaging.Lanczos), 1912, 390)
		}
		renderHeaderTextNumber, _ := gg.LoadFontFace(font, 65)
		getMainBgRender.SetFontFace(renderHeaderTextNumber)
		// white glow render
		getMainBgRender.SetRGB(1, 1, 1)
		getMainBgRender.DrawStringAnchored(getLink, 2021, 430, 0.4, 0.4)
		// avatar
		AvatarWaiter.Wait()
		getMainBgRender.DrawImage(getAvatarFormat.Image(), 2321, 230)
		getMainBgRender.Fill()
		// render
		CardRender(getMainBgRender, phidata)
		// draw bottom
		_ = getMainBgRender.LoadFontFace(font, 40)
		getMainBgRender.SetColor(color.White)
		getMainBgRender.Fill()
		getMainBgRender.DrawString("Generated By Lucy (HiMoYoBOT) | Designed By Eastown | Data From Phigros Library Project", 10, 5480)
		_ = getMainBgRender.SavePNG(engine.DataFolder() + "save/" + strconv.Itoa(int(ctx.Event.UserID)) + ".png")
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image("file:///"+file.BOTPATH+"/"+engine.DataFolder()+"save/"+strconv.Itoa(int(ctx.Event.UserID))+".png"))
	})
	engine.OnRegex(`^[! ！/]pgr\sroll\s(\d+)`).SetBlock(true).Limit(ctxext.LimitByGroup).Handle(func(ctx *zero.Ctx) {
		var wg sync.WaitGroup
		var avatarWaitGroup sync.WaitGroup
		var dataWaiter sync.WaitGroup
		var getMainBgRender *gg.Context
		var getAvatarFormat *gg.Context
		var setGlobalStat = true
		var phidata []byte
		wg.Add(1)
		avatarWaitGroup.Add(1)
		dataWaiter.Add(1)
		// get Session From Database.
		data := GetUserInfoFromDatabase(ctx.Event.UserID)
		getDataSession := data.PhiSession
		if getDataSession == "" {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("由于Session特殊性，请前往 https://phi.lemonkoi.one 获取绑定码进行绑定"))
			return
		}
		userData := GetUserInfoFromDatabase(ctx.Event.UserID)
		getRoll := ctx.State["regex_matched"].([]string)[1]
		getRollInt, err := strconv.ParseInt(getRoll, 10, 64)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请求roll不合法"))
			return
		}
		if getRollInt > 40 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("限制查询数为小于40"))
			return
		}
		getOverFlowNumber := getRollInt - 19
		if getOverFlowNumber <= 0 {
			getOverFlowNumber = 0
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("好哦~正在帮你请求，请稍等一下啦w~大约需要1-2分钟"))
		// data handling.
		go func() {
			defer dataWaiter.Done()
			getFullLink := router + "/api/phi/bests?session=" + userData.PhiSession + "&overflow=" + strconv.Itoa(int(getOverFlowNumber))
			phidata, _ = web.GetData(getFullLink)
			if phidata == nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("目前 Unoffical Phigros Library 暂时无法工作 请过一段时候尝试"))
				setGlobalStat = false
				return
			}
			err = json.Unmarshal(phidata, &phigrosB19)
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("发生解析错误: \n", err))
				setGlobalStat = false
				return
			}
			if !phigrosB19.Status {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("w? 貌似出现了一些问题x\n", phigrosB19.Message))
				setGlobalStat = false
				return
			}
		}()
		go func() {
			defer wg.Done()
			getRawBackground, _ := gg.LoadImage(backgroundRender)
			getMainBgRender = gg.NewContextForImage(imaging.Resize(getRawBackground, 2750, int(5250+getOverFlowNumber*200), imaging.Lanczos))
		}()
		go func() {
			defer avatarWaitGroup.Done()
			// draw Avatar, avatar from qq.
			avatarByte, err := http.Get("https://q4.qlogo.cn/g?b=qq&nk=" + strconv.FormatInt(ctx.Event.UserID, 10) + "&s=640")
			if err != nil {
				// pic lost.
				return
			}
			// avatar
			avatarByteUni, _, _ := image.Decode(avatarByte.Body)
			showUserAvatar := imaging.Resize(avatarByteUni, 250, 250, imaging.Lanczos)
			getAvatarFormat = gg.NewContext(250, 250)
			getAvatarFormat.DrawRoundedRectangle(0, 0, 248, 248, 20)
			getAvatarFormat.Clip()
			getAvatarFormat.DrawImage(showUserAvatar, 0, 0)
			getAvatarFormat.Fill()
		}()
		wg.Wait()
		_ = getMainBgRender.LoadFontFace(font, 30)
		// header background
		drawTriAngle(getMainBgRender, a, 0, 166, 1324, 410)
		getMainBgRender.SetRGBA255(0, 0, 0, 160)
		getMainBgRender.Fill()
		drawTriAngle(getMainBgRender, a, 1318, 192, 1600, 350)
		getMainBgRender.SetRGBA255(0, 0, 0, 160)
		getMainBgRender.Fill()
		drawTriAngle(getMainBgRender, a, 1320, 164, 6, 414)
		getMainBgRender.SetColor(color.White)
		getMainBgRender.Fill()
		// header background end.
		// load icon with other userinfo.
		getMainBgRender.SetColor(color.White)
		logo, _ := gg.LoadPNG(icon)
		getImageLogo := imaging.Resize(logo, 290, 290, imaging.Lanczos)
		getMainBgRender.DrawImage(getImageLogo, 50, 216)
		fontface, _ := gg.LoadFontFace(font, 90)
		getMainBgRender.SetFontFace(fontface)
		getMainBgRender.DrawString("Phigros", 422, 336)
		getMainBgRender.DrawString("RankingScore查询", 422, 462)
		dataWaiter.Wait()
		if !setGlobalStat {
			return
		}
		// draw userinfo path
		renderHeaderText, _ := gg.LoadFontFace(font, 54)
		getMainBgRender.SetFontFace(renderHeaderText)
		// wait data until fine.
		getMainBgRender.DrawString("Player: "+phigrosB19.Content.PlayerID, 1490, 300)
		getMainBgRender.DrawString("RankingScore: "+strconv.FormatFloat(phigrosB19.Content.RankingScore, 'f', 3, 64), 1490, 380)
		getMainBgRender.DrawString("ChanllengeMode: ", 1490, 460) // +56
		getColor, getLink := GetUserChallengeMode(phigrosB19.Content.ChallengeModeRank)
		if getColor != "" {
			getColorLink := ChanllengeMode + getColor + ".png"
			getColorImage, _ := gg.LoadImage(getColorLink)
			getMainBgRender.DrawImage(imaging.Resize(getColorImage, 238, 130, imaging.Lanczos), 1912, 390)
		}
		renderHeaderTextNumber, _ := gg.LoadFontFace(font, 65)
		getMainBgRender.SetFontFace(renderHeaderTextNumber)
		// white glow render
		getMainBgRender.SetRGB(1, 1, 1)
		getMainBgRender.DrawStringAnchored(getLink, 2021, 430, 0.4, 0.4)
		avatarWaitGroup.Wait()
		getMainBgRender.DrawImage(getAvatarFormat.Image(), 2321, 230)
		getMainBgRender.Fill()
		// render
		CardRender(getMainBgRender, phidata)
		// draw bottom
		_ = getMainBgRender.LoadFontFace(font, 40)
		getMainBgRender.SetColor(color.White)
		getMainBgRender.Fill()
		getMainBgRender.DrawString("Generated By Lucy (HiMoYoBOT) | Designed By Eastown | Data From Phigros Library Project", 10, float64(5110+getOverFlowNumber*200+100))
		_ = getMainBgRender.SavePNG(engine.DataFolder() + "save/" + "roll" + strconv.Itoa(int(ctx.Event.UserID)) + ".png")
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Image("file:///"+file.BOTPATH+"/"+engine.DataFolder()+"save/"+"roll"+strconv.Itoa(int(ctx.Event.UserID))+".png"))
	})
	engine.OnRegex(`^[! ！/]pgr\ssearch\s(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		// search a song: router: /api/phi/search
		getParams := ctx.State["regex_matched"].([]string)[1]
		queryData, err := web.GetData((router + "/api/phi/search?params=" + url.QueryEscape(getParams)))
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
			return
		}
		toString := helper.BytesToString(queryData)
		/*
			{
			    "status": true,
			    "content": {
			        "song_name": "DESTRUCTION 3,2,1",
			        "song_ratio": 0.45454545454545453,
			        "song_id": "DESTRUCTION321.Normal1zervsBrokenNerdz"
			    }
			}
		*/
		if !gjson.Get(toString, "status").Bool() {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("无法查询到目前歌曲"))
			return
		}
		//	getSongName := gjson.Get("content.song_name").Str
		//	getSongRatio := gjson.Get(toString, "content.song_ratio").Str
		getSongID := gjson.Get(toString, "content.song_id").Str
		// query to /api/phi/song , get song details.
		querySongDetailsData, err := web.GetData(router + "/api/phi/song?songid=" + getSongID)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
			return
		}
		var getSongDetails QuerySongDetailsGenerator
		err = json.Unmarshal(querySongDetailsData, &getSongDetails)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
			return
		}
		var setATStats string
		if getSongDetails.Content.Info.ChartDetail.At.Charter != "" {
			setATStats = "AT: " + strconv.FormatFloat(getSongDetails.Content.Info.ChartDetail.At.Rating, 'f', 1, 64) + " - " + getSongDetails.Content.Info.ChartDetail.At.Charter
		}
		result := fmt.Sprintf("歌曲名: %s \n 作曲: %s \n 曲绘: %s \n\n歌曲详细Rating：\nEZ: %s - %s \nHD: %s - %s \nIN: %s - %s \n%s", getSongDetails.Content.Info.Songname, getSongDetails.Content.Info.Composer, getSongDetails.Content.Info.Illustrator, floatToString(getSongDetails.Content.Info.ChartDetail.EZ.Rating), getSongDetails.Content.Info.ChartDetail.EZ.Charter, floatToString(getSongDetails.Content.Info.ChartDetail.HD.Rating), getSongDetails.Content.Info.ChartDetail.HD.Charter, floatToString(getSongDetails.Content.Info.ChartDetail.In.Rating), getSongDetails.Content.Info.ChartDetail.In.Charter, setATStats)
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(result))
	})
	engine.OnRegex(`^[! ！/]pgr\srandom`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		data, err := web.GetData(router + "/api/phi/rand")
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
			return
		}
		dataToSting := helper.BytesToString(data)
		getSongName := gjson.Get(dataToSting, "content.songname").Str
		getSongDiff := gjson.Get(dataToSting, "content.level").Str
		getSongRating := gjson.Get(dataToSting, "content.rating").Str
		getSongComposer := gjson.Get(dataToSting, "content.composer").Str
		// cute reply list ^^Meow.
		randomList := []string{"来看看这个哦~或许可以呢xwx", "今天先试试这个吧^^", "~Have a try qwq"}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(randomList[rand.Intn(len(randomList))], "\n"+getSongName+" - "+getSongDiff+" "+getSongRating+" by "+getSongComposer))
	})
}

func floatToString(floatnum float64) string {
	return strconv.FormatFloat(floatnum, 'f', 1, 64)
}
