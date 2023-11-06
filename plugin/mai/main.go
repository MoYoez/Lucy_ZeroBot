package mai

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"math/rand"
	rand2 "math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/FloatTech/floatbox/binary"
	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/gg"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/img/text"
	"github.com/PuerkitoBio/goquery"
	"github.com/fumiama/jieba/util/helper"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var (
	engine = control.Register("maidx", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "Hi NekoPachi!\n说明书: https://lucy.impart.icu",
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
		dataPlayer, err := QueryMaiBotDataFromQQ(int(uid))
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
		_, err := GetDefaultPlate(getDefaultInfo)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("设定的预设不正确"))
			return
		}
		_ = FormatUserDataBase(ctx.Event.UserID, GetUserInfoFromDatabase(ctx.Event.UserID), getDefaultInfo).BindUserDataBase()
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("已经设定好了哦~"))

	})
	//
	engine.OnFullMatch("mai什么").SetBlock(true).Handle(func(ctx *zero.Ctx) {
		// on load
		// get file first
		SongData, err := os.ReadFile(engine.DataFolder() + "music_data")
		if err != nil {
			panic(err)
		}
		var SongDataRandomTools MaiSongData
		err = json.Unmarshal(SongData, &SongDataRandomTools)
		if err != nil {
			ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("ERR: ", err))
			return
		}
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("来试试~ "+SongDataRandomTools[rand.Intn(len(SongDataRandomTools))].BasicInfo.Title+" 吧~"))
	})
	// TODO: 查歌.
	engine.OnRegex(`^[! ！/](mai|b50)\sunlock\s(.*)$`).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getDefaultInfo := ctx.State["regex_matched"].([]string)[2]
		UnlockerMai(getDefaultInfo, ctx)
	})
	engine.OnRegex(`^[! ！/](mai|b50)\sqrunlock`, QRHandler).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getPic := ctx.State["image_url"].([]string)[0]
		imageData, err := web.GetData(getPic)
		if err != nil {
			fmt.Print("get data err")
			return
		}
		// dec code.
		base64Image := base64.StdEncoding.EncodeToString(imageData)

		// online.

		reNewRequester, err := http.NewRequest("POST", "https://api.2weima.com/api/qrdecode", bytes.NewBuffer([]byte(`{"qr_base64": `+base64Image+`}`)))
		reNewRequester.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		reNewRequester.Header.Set("Authorization", "Bearer 4258|T2DpmaY2aqYFfhC43DAeK396JritJWHDR6U5WIcm") // 我知道你会看 但是这东西给你也没用.
		// GET CALLBACK
		client := &http.Client{}
		resp, err := client.Do(reNewRequester)
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}
		defer resp.Body.Close()
		datas, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		data := gjson.Get(helper.BytesToString(datas), "qr_content").String()
		UnlockerMai(data, ctx)
	})
}

func UnlockerMai(getDefaultInfo string, ctx *zero.Ctx) {
	htmlContent, err := http.Get("http://m1d1.xyz/")
	if err != nil {
		fmt.Println("无法访问网站:", err)
		return
	}
	defer htmlContent.Body.Close()

	cookies := htmlContent.Header["Set-Cookie"]
	htmlRawTo, err := io.ReadAll(htmlContent.Body)
	if err != nil {
		fmt.Println("无法访问网站:", err)
		return
	}
	htmlContentToStr := helper.BytesToString(htmlRawTo)
	if err != nil {
		ctx.SendChain(message.Text("解析Value出现问题，请稍后再尝试"))
		return
	}
	reader := strings.NewReader(htmlContentToStr)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		ctx.SendChain(message.Text("解析Value出现问题，请稍后再尝试"))
		return
	}

	inputValue := doc.Find("input[name=csrfmiddlewaretoken]").AttrOr("value", "")
	// get cookie

	// build post.
	form := url.Values{}
	form.Add("q", getDefaultInfo)
	form.Add("csrfmiddlewaretoken", inputValue)
	postData := form.Encode()
	req, err := http.NewRequest("POST", "http://m1d1.xyz/res/", strings.NewReader(postData))
	cookierSets := ""
	for i, cookier := range cookies {
		if i-1 == len(cookies) {
			cookierSets += cookier
		}
		cookierSets += cookier + ";"
	}
	fmt.Print(cookierSets)
	req.Header.Set("Set-Cookie", cookierSets)
	req.Header.Set("Cross-Origin-Opener-Policy", "same-origin")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		ctx.SendChain(message.Text("解析Value出现问题，请稍后再尝试"))
		return
	}
	// GET CALLBACK
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	docBack, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		panic(err)
	}
	pContent := docBack.Find("p").Text()
	ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(pContent))
}

func QRHandler(ctx *zero.Ctx) bool {
	if zero.HasPicture(ctx) {
		return true
	}
	// 没有图片就索取
	ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请提供二维码"))
	next := zero.NewFutureEvent("message", 999, true, ctx.CheckSession(), zero.HasPicture).Next()
	select {
	case <-time.After(time.Second * 30):
		return false
	case newCtx := <-next:
		ctx.State["image_url"] = newCtx.State["image_url"]
		ctx.Event.MessageID = newCtx.Event.MessageID
		return true
	}
}
