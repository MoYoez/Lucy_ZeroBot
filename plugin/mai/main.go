package mai

import (
	"bytes"
	"encoding/json"
	"image"
	rand2 "math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/gg"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/img/text"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

var (
	engine = control.Register("maidx", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "Hi NekoPachi!\n",
		PrivateDataFolder: "maidx",
	})
)

type MaiSongData []struct {
	Id     string    `json:"id"`
	Title  string    `json:"title"`
	Type   string    `json:"type"`
	Ds     []float64 `json:"ds"`
	Level  []string  `json:"level"`
	Cids   []int     `json:"cids"`
	Charts []struct {
		Notes   []int  `json:"notes"`
		Charter string `json:"charter"`
	} `json:"charts"`
	BasicInfo struct {
		Title       string `json:"title"`
		Artist      string `json:"artist"`
		Genre       string `json:"genre"`
		Bpm         int    `json:"bpm"`
		ReleaseDate string `json:"release_date"`
		From        string `json:"from"`
		IsNew       bool   `json:"is_new"`
	} `json:"basic_info"`
}

func init() {
	engine.OnRegex(`^[！!]chun$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		uid := ctx.Event.UserID
		dataPlayer, err := QueryChunDataFromQQ(int(uid))
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
			return
		}
		txt := HandleChunDataByUsingText(dataPlayer)
		base64Font, err := text.RenderToBase64(txt, text.BoldFontFile, 1920, 45)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
			return
		}
		ctx.SendChain(message.Image("base64://" + binary.BytesToString(base64Font)))
	})
	engine.OnRegex(`^[! ！/](mai|b50)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		uid := ctx.Event.UserID
		if GetUserSwitcherInfoFromDatabase(uid) == true {
			// use lxns checker service.
			getUserData := RequestBasicDataFromLxns(uid)
			if getUserData.Code != 200 {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("aw 出现了一点小错误~：\n - 请检查你是否有上传过数据并且绑定了QQ号\n - 请检查你的设置是否允许了第三方查看"))
				return
			}
			getGameUserData := RequestB50DataByFriendCode(getUserData.Data.FriendCode)
			if getGameUserData.Code != 200 {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("aw 出现了一点小错误~：\n - 请检查你是否有上传过数据并且绑定了QQ号\n - 请检查你的设置是否允许了第三方查看"))
				return
			}
			getImager, _ := ReFullPageRender(getGameUserData, getUserData, ctx)
			_ = gg.NewContextForImage(getImager).SavePNG(engine.DataFolder() + "save/" + "LXNS_" + strconv.Itoa(int(ctx.Event.UserID)) + ".png")
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Render User B50 : "+getUserData.Data.Name), message.Image(Saved+"LXNS_"+strconv.Itoa(int(ctx.Event.UserID))+".png"))
		} else {
			dataPlayer, err := QueryMaiBotDataFromQQ(int(uid), true)
			if err != nil {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
				return
			}
			var data player
			_ = json.Unmarshal(dataPlayer, &data)
			renderImg, plateStat := FullPageRender(data, ctx)
			tipPlate := ""
			getRand := rand2.Intn(10)
			if getRand == 8 {
				if !plateStat {
					tipPlate = "tips: 可以使用 ！mai plate xxx 来绑定称号~\n"
				}
			}
			_ = gg.NewContextForImage(renderImg).SavePNG(engine.DataFolder() + "save/" + strconv.Itoa(int(ctx.Event.UserID)) + ".png")
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Render User B50 : "+data.Username+"\n"+tipPlate), message.Image(Saved+strconv.Itoa(int(ctx.Event.UserID))+".png"))

		}
	})
	engine.OnRegex(`^[! ！/](b40)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		uid := ctx.Event.UserID
		dataPlayer, err := QueryMaiBotDataFromQQ(int(uid), false)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
			return
		}
		var data player
		_ = json.Unmarshal(dataPlayer, &data)
		renderImg, plateStat := FullPageRender(data, ctx)
		tipPlate := ""
		getRand := rand2.Intn(10)
		if getRand == 8 {
			if !plateStat {
				tipPlate = "tips: 可以使用 ！mai plate xxx 来绑定称号~\n"
			}
		}
		_ = gg.NewContextForImage(renderImg).SavePNG(engine.DataFolder() + "save/b40_" + strconv.Itoa(int(ctx.Event.UserID)) + ".png")
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Render User B40 : "+data.Username+"\n"+tipPlate), message.Image(Saved+"b40_"+strconv.Itoa(int(ctx.Event.UserID))+".png"))
	})
	engine.OnRegex(`^[! ！/](mai|b50)\s(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		matched := ctx.State["regex_matched"].([]string)[2]
		dataPlayer, err := QueryMaiBotDataFromUserName(matched)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
			return
		}
		var data player
		_ = json.Unmarshal(dataPlayer, &data)
		renderImg, plateStat := FullPageRender(data, ctx)
		tipPlate := ""
		if !plateStat {
			tipPlate = "tips: 可以使用 ！mai plate xxx 来绑定称号~"
		}
		_ = gg.NewContextForImage(renderImg).SavePNG(engine.DataFolder() + "save/" + strconv.Itoa(int(ctx.Event.UserID)) + ".png")
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Render User B50 : "+data.Username+"\n"+tipPlate+"\n"), message.Image(Saved+strconv.Itoa(int(ctx.Event.UserID))+".png"))
	})
	engine.OnRegex(`^[! ！/](mai|b50)\sswitch$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getBool := GetUserSwitcherInfoFromDatabase(ctx.Event.UserID)
		err := FormatUserSwitcher(ctx.Event.UserID, !getBool).ChangeUserSwitchInfoFromDataBase()
		if err != nil {
			panic(err)
		}
		var getEventText string
		// due to it changed, so reverse.
		if getBool == false {
			getEventText = "Lxns查分"
		} else {
			getEventText = "Diving Fish查分"
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经修改为"+getEventText))
	})
	engine.OnRegex(`^[! ！/](mai|b50)\splate\s(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getPlateInfo := ctx.State["regex_matched"].([]string)[2]
		_ = FormatUserDataBase(ctx.Event.UserID, getPlateInfo, GetUserDefaultinfoFromDatabase(ctx.Event.UserID)).BindUserDataBase()
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经将称号绑定上去了哦w"))
	})
	engine.OnRegex(`^[! ！/](mai|b50)\supload`, PictureHandler).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getPic := ctx.State["image_url"].([]string)[0]
		imageData, err := web.GetData(getPic)
		if err != nil {
			return
		}
		getRaw, _, err := image.Decode(bytes.NewReader(imageData))
		if err != nil {
			panic(err)
		}
		// pic Handler
		getRenderPlatePicRaw := gg.NewContext(1260, 210)
		getRenderPlatePicRaw.DrawRoundedRectangle(0, 0, 1260, 210, 10)
		getRenderPlatePicRaw.Clip()
		getHeight := getRaw.Bounds().Dy()
		getLength := getRaw.Bounds().Dx()
		var getHeightHandler, getLengthHandler int
		switch {
		case getHeight < 210 && getLength < 1260:
			getRaw = Resize(getRaw, 1260, 210)
			getHeightHandler = 0
			getLengthHandler = 0
		case getHeight < 210:
			getRaw = Resize(getRaw, getLength, 210)
			getHeightHandler = 0
			getLengthHandler = (getRaw.Bounds().Dx() - 1260) / 3 * -1
		case getLength < 1260:
			getRaw = Resize(getRaw, 1260, getHeight)
			getHeightHandler = (getRaw.Bounds().Dy() - 210) / 3 * -1
			getLengthHandler = 0
		default:
			getLengthHandler = (getRaw.Bounds().Dx() - 1260) / 3 * -1
			getHeightHandler = (getRaw.Bounds().Dy() - 210) / 3 * -1
		}
		getRenderPlatePicRaw.DrawImage(getRaw, getLengthHandler, getHeightHandler)
		getRenderPlatePicRaw.Fill()
		_ = getRenderPlatePicRaw.SavePNG(userPlate + strconv.Itoa(int(ctx.Event.UserID)) + ".png")
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经存入了哦w"))
	})
	engine.OnRegex(`^[! ！/](mai|b50)\sremove`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		_ = os.Remove(userPlate + strconv.Itoa(int(ctx.Event.UserID)) + ".png")
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经删掉了哦w"))
	})
	engine.OnRegex(`^[! ！/](mai|b50)\sdefault\splate\s(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getDefaultInfo := ctx.State["regex_matched"].([]string)[2]
		if getDefaultInfo == "" {
			_ = FormatUserDataBase(ctx.Event.UserID, GetUserInfoFromDatabase(ctx.Event.UserID), getDefaultInfo).BindUserDataBase()
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经恢复了正常~"))
			return
		}
		_, err := GetDefaultPlate(getDefaultInfo)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("设定的预设不正确"))
			return
		}
		_ = FormatUserDataBase(ctx.Event.UserID, GetUserInfoFromDatabase(ctx.Event.UserID), getDefaultInfo).BindUserDataBase()
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经设定好了哦~"))

	})
	engine.OnRegex(`^[! ！/](mai|b50)\sbind\s(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getDefaultInfo := ctx.State["regex_matched"].([]string)[2]
		indexReply := DecHashToRaw(getDefaultInfo)
		// get session.
		if indexReply == "" {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请前往 https://mai.lemonkoi.one 获取绑定码进行绑定"))
			return
		}
		getQID, getSessionID := RawJsonParse(indexReply)
		if getQID != ctx.Event.UserID {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请求Hash中QQ号不一致，请使用自己的号重新申请"))
			return
		}
		// check id
		getID := GetWahlapUserID(getSessionID)
		if getID == -1 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ID 无效或者是过期 ，请使用新的ID或者再次尝试"))
			return
		}
		// login.
		err := FormatUserIDDatabase(ctx.Event.UserID, strconv.Itoa(int(getID))).BindUserIDDataBase()
		if err != nil {
			panic(err)
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("绑定成功~"))
	})
	engine.OnRegex(`^[! ！/](mai|b50)\sunbind$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("取消绑定成功~"))
		RemoveUserIdFromDatabase(ctx.Event.UserID)
	})
	engine.OnRegex(`^[! ！/](mai|b50)\sunlock$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getMaiID := GetUserIDFromDatabase(ctx.Event.UserID)
		if getMaiID.Userid == "" {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请前往 https://mai.lemonkoi.one 获取绑定码进行绑定"))
			return
		}
		getCodeRaw, err := strconv.ParseInt(getMaiID.Userid, 10, 64)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Wahlap ServerERR: "+err.Error()))
			return
		}
		getCodeStat := Logout(getCodeRaw)
		if strings.Contains(getCodeStat, "{") == false {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("返回了错误.png, ERROR:"+getCodeStat))
			return
		}
		getCode := gjson.Get(getCodeStat, "returnCode").Int()
		if getCode == 1 {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("发信成功，服务器返回正常, 如果未生效请重新尝试"))
		} else {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("发信成功，但是服务器返回代码异常。"))
		}
	})
	engine.OnRegex(`^[! ！/](mai|b50)\stokenbind\s(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getDefaultInfo := ctx.State["regex_matched"].([]string)[2]
		FormatUserToken(strconv.FormatInt(ctx.Event.UserID, 10), getDefaultInfo).BindUserToken()
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("绑定成功~"))
	})
	engine.OnRegex(`^[! ！/](mai|b50)\supdate$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getID := ctx.Event.UserID
		getMaiID := GetUserIDFromDatabase(getID)
		if getMaiID.Userid == "" {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请前往 https://mai.lemonkoi.one 获取绑定码进行绑定"))
			return
		}
		getTokenId := GetUserToken(strconv.FormatInt(getID, 10))
		if getTokenId == "" {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请使用 /mai tokenbind <tokenid> 绑定水鱼查分器，其中 TokenID 从 https://www.diving-fish.com/maimaidx/prober 用户设置中拿到"))
			return
		}
		if !CheckTheTicketIsValid(getTokenId) {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("此 Token 不合法 ，请重新绑定"))
			return
		}
		// token is valid, get data.
		getIntID, _ := strconv.ParseInt(getMaiID.Userid, 10, 64)
		getFullData := GetMusicList(getIntID, 0, 1000)
		if strings.Contains(getFullData, "{") == false {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("返回了错误.png, ERROR:"+getFullData))
			return
		}
		var jsonMashell UserMusicListStruct
		err := json.Unmarshal(helper.StringToBytes(getFullData), &jsonMashell)
		if err != nil {
			panic(err)
		}
		getFullDataStruct := convert(jsonMashell)
		jsonDumper := getFullDataStruct
		jsonDumperFull, err := json.Marshal(jsonDumper)
		if err != nil {
			panic(err)
		}
		// upload to diving fish api
		req, err := http.NewRequest("POST", "https://www.diving-fish.com/api/maimaidxprober/player/update_records", bytes.NewBuffer(jsonDumperFull))
		if err != nil {
			// Handle error
			panic(err)
		}
		req.Header.Set("Import-Token", getTokenId)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		//	NewReader, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("Update CODE:"+strconv.Itoa(resp.StatusCode)))
	})
	engine.OnRegex(`^[! ！/](mai|b50)\sregion$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getID := ctx.Event.UserID
		getMaiID := GetUserIDFromDatabase(getID)
		if getMaiID.Userid == "" {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请前往 https://mai.lemonkoi.one 获取绑定码进行绑定"))
			return
		}
		getCodeRaw, err := strconv.ParseInt(getMaiID.Userid, 10, 64)
		if err != nil {
			panic(err)
		}
		getReplyMsg := GetUserRegion(getCodeRaw)
		if strings.Contains(getReplyMsg, "{") == false {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("返回了错误.png, ERROR:"+getReplyMsg))
			return
		}
		var MixedMagic GetUserRegionStruct
		json.Unmarshal(helper.StringToBytes(getReplyMsg), &MixedMagic)
		var returnText string
		for _, onlistLoader := range MixedMagic.UserRegionList {
			returnText = returnText + MixedRegionWriter(onlistLoader.RegionId-1, onlistLoader.PlayCount, onlistLoader.Created) + "\n\n"
		}
		if returnText == "" {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("目前 Lucy 没有查到您的游玩记录哦~"))
			return
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("目前查询到您的游玩记录如下: \n\n"+returnText))
	})
	engine.OnRegex(`^(网咋样|[! ！/](mai|b50)\sstatus$)`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		// getWebStatus
		//	getWebStatus := ReturnWebStatus()
		getZlibError := ReturnZlibError()
		getPlayedStatus, err := web.GetData("https://maihook.lemonkoi.one/api/calc")
		if err != nil {
			return
		}
		var playerStatus RealConvertPlay
		json.Unmarshal(getPlayedStatus, &playerStatus)
		// 20s one request.
		var getLucyRespHandler int
		if getZlibError.Full.Field3 < 180 {
			getLucyRespHandler = getZlibError.Full.Field3
		} else {
			getLucyRespHandler = getZlibError.Full.Field3 - 180
		}
		getLucyRespHandlerStr := strconv.Itoa(getLucyRespHandler)

		getZlibWord := "Zlib 压缩跳过率: \n" + "10mins (" + ConvertZlib(getZlibError.ZlibError.Field1, getZlibError.Full.Field1) + " Loss)\n" + "30mins (" + ConvertZlib(getZlibError.ZlibError.Field2, getZlibError.Full.Field2) + " Loss)\n" + "60mins (" + ConvertZlib(getZlibError.ZlibError.Field3, getZlibError.Full.Field3) + " Loss)\n"
		getRealStatus := "\n以下数据来源于mai机台的数据反馈\n"
		//	getWebStatusCount := "Web Uptime Ping:\n * MaimaiDXCN: " + ConvertFloat(getWebStatus.Details.MaimaiDXCN.Uptime*100) + "%\n * MaimaiDXCN Main Server: " + ConvertFloat(getWebStatus.Details.MaimaiDXCNMain.Uptime*100) + "%\n * MaimaiDXCN Title Server: " + ConvertFloat(float64(getWebStatus.Details.MaimaiDXCNTitle.Uptime*100)) + "%\n * MaimaiDXCN Update Server: " + ConvertFloat(float64(getWebStatus.Details.MaimaiDXCNUpdate.Uptime*100)) + "%\n * MaimaiDXCN NetLogin Server: " + ConvertFloat(getWebStatus.Details.MaimaiDXCNNetLogin.Uptime*100) + "%\n * MaimaiDXCN Net Server: " + ConvertFloat(getWebStatus.Details.MaimaiDXCNDXNet.Uptime*100) + "%\n"
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("* Zlib 压缩跳过率可以很好的反馈当前 MaiNet (Wahlap Service) 当前负载的情况，根据样本 + Lucy处理情况 来判断 \n* 错误率收集则来源于 cici 分布式部署监测工具，反应各地区真实mai游玩错误情况 \n* 在 1小时 内，Lucy 共处理了 "+getLucyRespHandlerStr+"次 请求💫，其中详细数据如下:\n\n"+getZlibWord+getRealStatus+"\n"+ConvertRealPlayWords(playerStatus)+"\n* Zlib 3% Loss 以下则 基本上可以正常游玩\n* 10% Loss 则会有明显断网现象(请准备小黑屋工具)\n* 30% Loss 则无法正常游玩(即使使用小黑屋工具) "))
	})

}
