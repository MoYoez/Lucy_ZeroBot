package mai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/FloatTech/floatbox/file"
	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/gg"
	"github.com/FloatTech/imgfactory"
	"github.com/tidwall/gjson"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/text/width"
)

type RealConvertPlay struct {
	ReturnValue []struct {
		SkippedCount  int `json:"skippedCount"`
		RetriedCount  int `json:"retriedCount"`
		RetryCountSum int `json:"retryCountSum"`
		TotalCount    int `json:"totalCount"`
		FailedCount   int `json:"failedCount"`
	} `json:"returnValue"`
}

type WebPingStauts struct {
	Details struct {
		MaimaiDXCN struct {
			Uptime float64 `json:"uptime"`
		} `json:"maimai DX CN"`
		MaimaiDXCNDXNet struct {
			Uptime float64 `json:"uptime"`
		} `json:"maimai DX CN DXNet"`
		MaimaiDXCNMain struct {
			Uptime float64 `json:"uptime"`
		} `json:"maimai DX CN Main"`
		MaimaiDXCNNetLogin struct {
			Uptime float64 `json:"uptime"`
		} `json:"maimai DX CN NetLogin"`
		MaimaiDXCNTitle struct {
			Uptime float64 `json:"uptime"`
		} `json:"maimai DX CN Title"`
		MaimaiDXCNUpdate struct {
			Uptime float64 `json:"uptime"`
		} `json:"maimai DX CN Update"`
	} `json:"details"`
	Status bool `json:"status"`
}

type ZlibErrorStatus struct {
	Full struct {
		Field1 int `json:"10"`
		Field2 int `json:"30"`
		Field3 int `json:"60"`
	} `json:"full"`
	FullError struct {
		Field1 int `json:"10"`
		Field2 int `json:"30"`
		Field3 int `json:"60"`
	} `json:"full_Error"`
	ZlibError struct {
		Field1 int `json:"10"`
		Field2 int `json:"30"`
		Field3 int `json:"60"`
	} `json:"zlib_Error"`
}

type player struct {
	AdditionalRating int `json:"additional_rating"`
	Charts           struct {
		Dx []struct {
			Achievements float64 `json:"achievements"`
			Ds           float64 `json:"ds"`
			DxScore      int     `json:"dxScore"`
			Fc           string  `json:"fc"`
			Fs           string  `json:"fs"`
			Level        string  `json:"level"`
			LevelIndex   int     `json:"level_index"`
			LevelLabel   string  `json:"level_label"`
			Ra           int     `json:"ra"`
			Rate         string  `json:"rate"`
			SongId       int     `json:"song_id"`
			Title        string  `json:"title"`
			Type         string  `json:"type"`
		} `json:"dx"`
		Sd []struct {
			Achievements float64 `json:"achievements"`
			Ds           float64 `json:"ds"`
			DxScore      int     `json:"dxScore"`
			Fc           string  `json:"fc"`
			Fs           string  `json:"fs"`
			Level        string  `json:"level"`
			LevelIndex   int     `json:"level_index"`
			LevelLabel   string  `json:"level_label"`
			Ra           int     `json:"ra"`
			Rate         string  `json:"rate"`
			SongId       int     `json:"song_id"`
			Title        string  `json:"title"`
			Type         string  `json:"type"`
		} `json:"sd"`
	} `json:"charts"`
	Nickname string      `json:"nickname"`
	Plate    string      `json:"plate"`
	Rating   int         `json:"rating"`
	UserData interface{} `json:"user_data"`
	UserId   interface{} `json:"user_id"`
	Username string      `json:"username"`
}

type playerData struct {
	Achievements float64 `json:"achievements"`
	Ds           float64 `json:"ds"`
	DxScore      int     `json:"dxScore"`
	Fc           string  `json:"fc"`
	Fs           string  `json:"fs"`
	Level        string  `json:"level"`
	LevelIndex   int     `json:"level_index"`
	LevelLabel   string  `json:"level_label"`
	Ra           int     `json:"ra"`
	Rate         string  `json:"rate"`
	SongId       int     `json:"song_id"`
	Title        string  `json:"title"`
	Type         string  `json:"type"`
}

type chun struct {
	Nickname string  `json:"nickname"`
	Rating   float64 `json:"rating"`
	Records  struct {
		B30 []struct {
			Cid        int     `json:"cid"`
			Ds         float64 `json:"ds"`
			Fc         string  `json:"fc"`
			Level      string  `json:"level"`
			LevelIndex int     `json:"level_index"`
			LevelLabel string  `json:"level_label"`
			Mid        int     `json:"mid"`
			Ra         float64 `json:"ra"`
			Score      int     `json:"score"`
			Title      string  `json:"title"`
		} `json:"b30"`
		R10 []struct {
			Cid        int     `json:"cid"`
			Ds         float64 `json:"ds"`
			Fc         string  `json:"fc"`
			Level      string  `json:"level"`
			LevelIndex int     `json:"level_index"`
			LevelLabel string  `json:"level_label"`
			Mid        int     `json:"mid"`
			Ra         float64 `json:"ra"`
			Score      int     `json:"score"`
			Title      string  `json:"title"`
		} `json:"r10"`
	} `json:"records"`
	Username string `json:"username"`
}

var (
	loadMaiPic        = Root + "pic/"
	defaultCoverLink  = Root + "default_cover.png"
	typeImageDX       = loadMaiPic + "chart_type_dx.png"
	typeImageSD       = loadMaiPic + "chart_type_sd.png"
	titleFontPath     = maifont + "NotoSansSC-Bold.otf"
	UniFontPath       = maifont + "Montserrat-Bold.ttf"
	nameFont          = maifont + "NotoSansSC-Regular.otf"
	maifont           = Root + "font/"
	b50bgOriginal     = loadMaiPic + "b50_bg.png"
	b50bg             = loadMaiPic + "b50_bg.png"
	b50Custom         = loadMaiPic + "b50_bg_custom.png"
	Root              = engine.DataFolder() + "resources/maimai/"
	userPlate         = engine.DataFolder() + "user/"
	Saved             = "file:///" + file.BOTPATH + "/" + engine.DataFolder() + "save/"
	titleFont         font.Face
	scoreFont         font.Face
	rankFont          font.Face
	levelFont         font.Face
	ratingFont        font.Face
	nameTypeFont      font.Face
	diffColor         []color.RGBA
	ratingBgFilenames = []string{
		"rating_white.png",
		"rating_blue.png",
		"rating_green.png",
		"rating_yellow.png",
		"rating_red.png",
		"rating_purple.png",
		"rating_copper.png",
		"rating_silver.png",
		"rating_gold.png",
		"rating_rainbow.png",
	}
)

func HandleChunDataByUsingText(handleJson []byte) string {
	var chunData chun
	_ = json.Unmarshal(handleJson, &chunData)

	/*
					Chun | Player: MoYoez (MoeMagicMango) | Rating: 10.442
			-- b30
				- 1. Lemegeton -little key of solomon- (Expert 13)  FC?  | Score: 974529 | rating: 13.17

			-- r10
				- 1. ....
		Generated By Lucy,code with lazy.
	*/
	var wg sync.WaitGroup
	wg.Add(2)
	getUserName := chunData.Username
	getNickName := chunData.Nickname
	getRating := strconv.FormatFloat(chunData.Rating, 'f', 3, 64)
	getB30Length := len(chunData.Records.B30)
	getR10Length := len(chunData.Records.R10)
	generateInitHeader := "Chun | Player: " + getNickName + "( " + getUserName + " ) | Rating: " + getRating
	var getB30Str string
	var getR10Str string
	go func() {
		defer wg.Done()
		for i := 0; i < getB30Length; i++ {
			getB30Str += strconv.Itoa(i+1) + ". " + chunData.Records.B30[i].Title + "( " + chunData.Records.B30[i].LevelLabel + " " + chunData.Records.B30[i].Level + " ) " + chunData.Records.B30[i].Fc + " | Score: " + strconv.Itoa(chunData.Records.B30[i].Score) + " | Rating: " + strconv.FormatFloat(chunData.Records.B30[i].Ra, 'f', 3, 64) + "\n"
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < getR10Length; i++ {
			getR10Str += strconv.Itoa(i+1) + ". " + chunData.Records.R10[i].Title + "( " + chunData.Records.R10[i].LevelLabel + " " + chunData.Records.R10[i].Level + " ) " + chunData.Records.R10[i].Fc + " | Score: " + strconv.Itoa(chunData.Records.R10[i].Score) + " | Rating: " + strconv.FormatFloat(chunData.Records.R10[i].Ra, 'f', 3, 64) + "\n"
		}
	}()
	wg.Wait()
	return generateInitHeader + "\n -- B30\n" + getB30Str + " -- R10 \n " + getR10Str + "\n Generated By Lucy,Code With Lazy."
}

func init() {
	if _, err := os.Stat(userPlate); os.IsNotExist(err) {
		err := os.MkdirAll(userPlate, 0777)
		if err != nil {
			panic(err)
		}
	}
	nameTypeFont = LoadFontFace(nameFont, 36)
	titleFont = LoadFontFace(titleFontPath, 20)
	scoreFont = LoadFontFace(UniFontPath, 32)
	rankFont = LoadFontFace(UniFontPath, 24)
	levelFont = LoadFontFace(UniFontPath, 20)
	ratingFont = LoadFontFace(UniFontPath, 24)
	diffColor = []color.RGBA{
		{69, 193, 36, 255},
		{255, 186, 1, 255},
		{255, 90, 102, 255},
		{134, 49, 200, 255},
		{207, 144, 240, 255},
	}
}

// FullPageRender  Render Full Page
func FullPageRender(data player, ctx *zero.Ctx) (raw image.Image, stat bool) {
	// muilt-threading.
	var avatarHandler sync.WaitGroup
	avatarHandler.Add(1)
	var getAvatarFormat *gg.Context
	// avatar handler.
	go func() {
		// avatar Round Style
		defer avatarHandler.Done()
		avatarByte, err := http.Get("https://q4.qlogo.cn/g?b=qq&nk=" + strconv.FormatInt(ctx.Event.UserID, 10) + "&s=640")
		if err != nil {
			return
		}
		avatarByteUni, _, _ := image.Decode(avatarByte.Body)
		avatarFormat := imgfactory.Size(avatarByteUni, 180, 180)
		getAvatarFormat = gg.NewContext(180, 180)
		getAvatarFormat.DrawRoundedRectangle(0, 0, 178, 178, 20)
		getAvatarFormat.Clip()
		getAvatarFormat.DrawImage(avatarFormat.Image(), 0, 0)
		getAvatarFormat.Fill()
	}()
	userPlatedCustom := GetUserDefaultinfoFromDatabase(ctx.Event.UserID)
	// render Header.
	b50Render := gg.NewContext(2090, 1660)
	rawPlateData, errs := gg.LoadImage(userPlate + strconv.Itoa(int(ctx.Event.UserID)) + ".png")
	if errs == nil {
		b50bg = b50Custom
		b50Render.DrawImage(rawPlateData, 595, 30)
		b50Render.Fill()
	} else {
		if userPlatedCustom != "" {
			b50bg = b50Custom
			images, _ := GetDefaultPlate(userPlatedCustom)
			b50Render.DrawImage(images, 595, 30)
			b50Render.Fill()
		} else {
			// show nil
			b50bg = b50bgOriginal
		}
	}
	getContent, _ := gg.LoadImage(b50bg)
	b50Render.DrawImage(getContent, 0, 0)
	b50Render.Fill()
	// render user info
	avatarHandler.Wait()
	b50Render.DrawImage(getAvatarFormat.Image(), 610, 50)
	b50Render.Fill()
	// render Userinfo
	b50Render.SetFontFace(nameTypeFont)
	b50Render.SetColor(color.Black)
	b50Render.DrawStringAnchored(width.Widen.String(data.Nickname), 825, 160, 0, 0)
	b50Render.Fill()
	b50Render.SetFontFace(titleFont)
	setPlateLocalStatus := GetUserInfoFromDatabase(ctx.Event.UserID)
	var dataPlate bool
	if setPlateLocalStatus != "" {
		data.Plate = setPlateLocalStatus
		dataPlate = true
	} else {
		dataPlate = false
	}
	b50Render.DrawStringAnchored(strings.Join(strings.Split(data.Plate, ""), " "), 1050, 207, 0.5, 0.5)
	b50Render.Fill()
	getRating := getRatingBg(data.Rating)
	getRatingBG, err := gg.LoadImage(loadMaiPic + getRating)
	if err != nil {
		panic(err)
	}
	b50Render.DrawImage(getRatingBG, 800, 40)
	b50Render.Fill()
	// render Rank
	imgs, err := GetRankPicRaw(data.AdditionalRating)
	if err != nil {
		panic(err)
	}
	b50Render.DrawImage(imgs, 1080, 50)
	b50Render.Fill()
	// draw number
	b50Render.SetFontFace(scoreFont)
	b50Render.SetRGBA255(236, 219, 113, 255)
	b50Render.DrawStringAnchored(strconv.Itoa(data.Rating), 1056, 60, 1, 1)
	b50Render.Fill()
	// Render Card Type
	getSDLength := len(data.Charts.Sd)
	getDXLength := len(data.Charts.Dx)
	getDXinitX := 45
	getDXinitY := 1225
	getInitX := 45
	getInitY := 285
	var i int
	for i = 0; i < getSDLength; i++ {
		b50Render.DrawImage(RenderCard(data.Charts.Sd[i], i+1, false), getInitX, getInitY)
		getInitX += 400
		if getInitX == 2045 {
			getInitX = 45
			getInitY += 125
		}
	}
	for dx := 0; dx < getDXLength; dx++ {
		b50Render.DrawImage(RenderCard(data.Charts.Dx[dx], dx+1, false), getDXinitX, getDXinitY)
		getDXinitX += 400
		if getDXinitX == 2045 {
			getDXinitX = 45
			getDXinitY += 125
		}
	}
	return b50Render.Image(), dataPlate
}

// RenderCard Main Lucy Render Page
func RenderCard(data playerData, num int, isSimpleRender bool) image.Image {
	getType := data.Type
	var CardBackGround string
	var multiTypeRender sync.WaitGroup
	var CoverDownloader sync.WaitGroup
	CoverDownloader.Add(1)
	multiTypeRender.Add(1)
	// choose Type.
	if getType == "SD" {
		CardBackGround = typeImageSD
	} else {
		CardBackGround = typeImageDX
	}
	charCount := 0.0
	setBreaker := false
	var truncated string
	var charFloatNum float64
	getSongName := data.Title
	var getSongId string
	switch {
	case data.SongId < 1000:
		getSongId = fmt.Sprintf("%05d", data.SongId)
	case data.SongId < 10000:
		getSongId = fmt.Sprintf("1%d", data.SongId)
	default:
		getSongId = strconv.Itoa(data.SongId)
	}
	var Image image.Image
	go func() {
		defer CoverDownloader.Done()
		Image, _ = GetCover(getSongId)
	}()
	// set rune count
	go func() {
		defer multiTypeRender.Done()
		for _, runeValue := range getSongName {
			charWidth := utf8.RuneLen(runeValue)
			if charWidth == 3 {
				charFloatNum = 1.5
			} else {
				charFloatNum = float64(charWidth)
			}
			if charCount+charFloatNum > 19 {
				setBreaker = true
				break
			}
			truncated += string(runeValue)
			charCount += charFloatNum
		}
		if setBreaker {
			getSongName = truncated + ".."
		} else {
			getSongName = truncated
		}
	}()
	loadSongType, _ := gg.LoadImage(CardBackGround)
	// draw pic
	drawBackGround := gg.NewContextForImage(GetChartType(data.LevelLabel))
	// draw song pic
	CoverDownloader.Wait()
	drawBackGround.DrawImage(Image, 25, 25)
	// draw name
	drawBackGround.SetColor(color.White)
	drawBackGround.SetFontFace(titleFont)
	multiTypeRender.Wait()
	drawBackGround.DrawStringAnchored(getSongName, 130, 32.5, 0, 0.5)
	drawBackGround.Fill()
	// draw acc
	drawBackGround.SetFontFace(scoreFont)
	drawBackGround.DrawStringAnchored(strconv.FormatFloat(data.Achievements, 'f', 4, 64)+"%", 129, 62.5, 0, 0.5)
	// draw rate
	drawBackGround.DrawImage(GetRateStatusAndRenderToImage(data.Rate), 305, 45)
	drawBackGround.Fill()
	drawBackGround.SetFontFace(rankFont)
	drawBackGround.SetColor(diffColor[data.LevelIndex])
	if !isSimpleRender {
		drawBackGround.DrawString("#"+strconv.Itoa(num), 130, 111)
	}
	drawBackGround.FillPreserve()
	// draw rest of card.
	drawBackGround.SetFontFace(levelFont)
	drawBackGround.DrawString(strconv.FormatFloat(data.Ds, 'f', -1, 64), 195, 111)
	drawBackGround.FillPreserve()
	drawBackGround.SetFontFace(ratingFont)
	drawBackGround.DrawString("▶", 235, 111)
	drawBackGround.FillPreserve()
	drawBackGround.SetFontFace(ratingFont)
	drawBackGround.DrawString(strconv.Itoa(data.Ra), 250, 111)
	drawBackGround.FillPreserve()
	if data.Fc != "" {
		drawBackGround.DrawImage(LoadComboImage(data.Fc), 290, 84)
	}
	if data.Fs != "" {
		drawBackGround.DrawImage(LoadSyncImage(data.Fs), 325, 84)
	}
	drawBackGround.DrawImage(loadSongType, 68, 88)
	return drawBackGround.Image()
}

func GetRankPicRaw(id int) (image.Image, error) {
	var idStr string
	if id < 10 {
		idStr = "0" + strconv.FormatInt(int64(id), 10)
	} else {
		idStr = strconv.FormatInt(int64(id), 10)
	}
	if id == 22 {
		idStr = "21"
	}
	data := Root + "rank/UI_CMN_DaniPlate_" + idStr + ".png"
	imgRaw, err := gg.LoadImage(data)
	if err != nil {
		return nil, err
	}
	return imgRaw, nil
}

func GetDefaultPlate(id string) (image.Image, error) {
	data := Root + "plate/plate_" + id + ".png"
	imgRaw, err := gg.LoadImage(data)
	if err != nil {
		return nil, err
	}
	return imgRaw, nil
}

// GetCover Careful The nil data
func GetCover(id string) (image.Image, error) {
	fileName := id + ".png"
	filePath := Root + "cover/" + fileName
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Auto download cover from diving fish's site
		downloadURL := "https://www.diving-fish.com/covers/" + fileName
		cover, err := downloadImage(downloadURL)
		if err != nil {
			return LoadPictureWithResize(defaultCoverLink, 90, 90), nil
		}
		saveImage(cover, filePath)
	}
	imageFile, err := os.Open(filePath)
	if err != nil {
		return LoadPictureWithResize(defaultCoverLink, 90, 90), nil
	}
	defer func(imageFile *os.File) {
		err := imageFile.Close()
		if err != nil {
			panic(err)
		}
	}(imageFile)
	img, _, err := image.Decode(imageFile)
	if err != nil {
		return LoadPictureWithResize(defaultCoverLink, 90, 90), nil
	}
	return Resize(img, 90, 90), nil
}

// Resize Image width height
func Resize(image image.Image, w int, h int) image.Image {
	return imgfactory.Size(image, w, h).Image()
}

// LoadPictureWithResize Load Picture
func LoadPictureWithResize(link string, w int, h int) image.Image {
	getImage, err := gg.LoadImage(link)
	if err != nil {
		return nil
	}
	return Resize(getImage, w, h)
}

// GetRateStatusAndRenderToImage Get Rate
func GetRateStatusAndRenderToImage(rank string) image.Image {
	// Load rank images
	return LoadPictureWithResize(loadMaiPic+"rate_"+rank+".png", 80, 40)
}

// GetChartType Get Chart Type
func GetChartType(chart string) image.Image {
	data, _ := gg.LoadImage(loadMaiPic + "chart_" + NoHeadLineCase(chart) + ".png")
	return data
}

// LoadComboImage Load combo images
func LoadComboImage(imageName string) image.Image {
	link := loadMaiPic + "combo_" + imageName + ".png"
	return LoadPictureWithResize(link, 60, 40)
}

// LoadSyncImage Load sync images
func LoadSyncImage(imageName string) image.Image {
	link := loadMaiPic + "sync_" + imageName + ".png"
	return LoadPictureWithResize(link, 60, 40)
}

// NoHeadLineCase No HeadLine.
func NoHeadLineCase(word string) string {
	text := strings.ToLower(word)
	textNewer := strings.ReplaceAll(text, ":", "")
	return textNewer
}

// LoadFontFace load font face once before running, to work it quickly and save memory.
func LoadFontFace(filePath string, size float64) font.Face {
	fontFile, _ := os.ReadFile(filePath)
	fontFileParse, _ := opentype.Parse(fontFile)
	fontFace, _ := opentype.NewFace(fontFileParse, &opentype.FaceOptions{Size: size, DPI: 70, Hinting: font.HintingFull})
	return fontFace
}

// Inline Code.
func saveImage(img image.Image, path string) {
	files, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer func(files *os.File) {
		err := files.Close()
		if err != nil {
			panic(err)
		}
	}(files)
	err = png.Encode(files, img)
	if err != nil {
		log.Fatal(err)
	}
}

func downloadImage(url string) (image.Image, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(response.Body)
	img, _, err := image.Decode(response.Body)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func getRatingBg(rating int) string {
	index := 0
	switch {
	case rating >= 15000:
		index++
		fallthrough
	case rating >= 14000:
		index++
		fallthrough
	case rating >= 13000:
		index++
		fallthrough
	case rating >= 12000:
		index++
		fallthrough
	case rating >= 10000:
		index++
		fallthrough
	case rating >= 7000:
		index++
		fallthrough
	case rating >= 4000:
		index++
		fallthrough
	case rating >= 2000:
		index++
		fallthrough
	case rating >= 1000:
		index++
	}
	return ratingBgFilenames[index]
}

// PictureHandler MaiPictureHandler Handler Mai Pic
func PictureHandler(ctx *zero.Ctx) bool {
	if zero.HasPicture(ctx) {
		return true
	}
	// 没有图片就索取
	ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text("请提供一张图片，图片大小比例适应为6:1 (1260x210) ,如果图片不适应将会自动剪辑到合适大小"))
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

func CheckTheTicketIsValid(ticket string) bool {
	getData, err := web.GetData("https://www.diving-fish.com/api/maimaidxprober/token_available?token=" + ticket)
	if err != nil {
		panic(err)
	}
	result := gjson.Get(helper.BytesToString(getData), "message").String()
	if result == "ok" {
		return true
	}
	return false
}

func convert(listStruct UserMusicListStruct) []InnerStructChanger {
	getRequest, err := web.GetData("https://www.diving-fish.com/api/maimaidxprober/music_data")
	if err != nil {
		panic(err)
	}
	var divingfishMusicData []DivingFishMusicDataStruct
	err = json.Unmarshal(getRequest, &divingfishMusicData)
	if err != nil {
		panic(err)
	}
	mdMap := make(map[string]DivingFishMusicDataStruct)
	for _, m := range divingfishMusicData {
		mdMap[m.Id] = m
	}
	var dest []InnerStructChanger
	for _, musicList := range listStruct.UserMusicList {
		for _, musicDetailedList := range musicList.UserMusicDetailList {
			level := musicDetailedList.Level
			achievement := math.Min(1010000, float64(musicDetailedList.Achievement))
			fc := []string{"", "fc", "fcp", "ap", "app"}[musicDetailedList.ComboStatus]
			fs := []string{"", "fs", "fsp", "fsd", "fsdp"}[musicDetailedList.SyncStatus]
			dxScore := musicDetailedList.DeluxscoreMax
			dest = append(dest, InnerStructChanger{
				Title:        mdMap[strconv.Itoa(musicDetailedList.MusicId)].Title,
				Type:         mdMap[strconv.Itoa(musicDetailedList.MusicId)].Type,
				LevelIndex:   level,
				Achievements: (achievement) / 10000,
				Fc:           fc,
				Fs:           fs,
				DxScore:      dxScore,
			})
		}
	}
	return dest
}

// MixedRegionWriter Some Mixed Magic, looking for your region information.
func MixedRegionWriter(regionID int, playCount int, createdDate string) string {
	getCountryID := returnCountryID(regionID)
	return fmt.Sprintf(" - 在 regionID 为 %d (%s) 的省/直辖市 游玩过 %d 次, 第一次游玩时间于 %s", regionID+1, getCountryID, playCount, createdDate)
}

// ReportToEndPoint Report Some Error To Wahlap Server.
func ReportToEndPoint(getReport int, getReportType string) string {
	url := "https://maihook.lemonkoi.one/api/zlib?report=" + strconv.Itoa(getReport) + "&reportType=" + getReportType
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Add("authkey", authKey)
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "maihook.lemonkoi.one")
	req.Header.Add("Connection", "keep-alive")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(body)
}

// ReturnZlibError Return Zlib ERROR
func ReturnZlibError() ZlibErrorStatus {
	url := "https://maihook.lemonkoi.one/api/zlib"
	method := "GET"
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return ZlibErrorStatus{}
	}
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Host", "maihook.lemonkoi.one")
	req.Header.Add("Connection", "keep-alive")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ZlibErrorStatus{}
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ZlibErrorStatus{}
	}
	var returnData ZlibErrorStatus
	json.Unmarshal(body, &returnData)
	return returnData
}

func ConvertZlib(value, total int) string {
	if total == 0 {
		return "0.000%"
	}
	percentage := float64(value) / float64(total) * 100
	return fmt.Sprintf("%.3f%%", percentage)
}

func ConvertRealPlayWords(retry RealConvertPlay) string {
	var pickedWords string
	var count = 0
	var header = " - 错误率数据收集自机台的真实网络通信，可以反映舞萌 DX 的网络状况。\n"

	for _, word := range retry.ReturnValue {
		var timeCount int
		var UserReturnLogs string
		switch {
		case count == 0:
			timeCount = 10
		case count == 1:
			timeCount = 30
		case count == 2:
			timeCount = 60
		}

		if word.TotalCount < 20 {
			UserReturnLogs = "没有收集到足够的数据进行分析~"
		} else {
			totalSuccess := word.TotalCount - word.FailedCount
			skippedRate := float64(word.SkippedCount) / float64(totalSuccess) * 100
			otherErrorRate := float64(word.RetryCountSum) / float64(totalSuccess+word.RetryCountSum) * 100
			overallErrorRate := (float64(word.SkippedCount+word.RetryCountSum) / float64(totalSuccess+word.RetryCountSum)) * 100
			skippedRate = math.Round(skippedRate*100) / 100
			otherErrorRate = math.Round(otherErrorRate*100) / 100
			overallErrorRate = math.Round(overallErrorRate*100) / 100
			UserReturnLogs = fmt.Sprintf("共 %d 个成功的请求中，有 %d 次未压缩（%.2f%%），有 %d 个请求共 %d 次其他错误（%.2f%%），整体错误率为 %.2f%%。", totalSuccess, word.SkippedCount, skippedRate, word.RetriedCount, word.RetryCountSum, otherErrorRate, overallErrorRate)
		}
		pickedWords = pickedWords + fmt.Sprintf("\n - 在 %d 分钟内%s", timeCount, UserReturnLogs)
		count = count + 1

	}
	return header + pickedWords + "\n"
}

func RequestReferSong(friendID int64, songID int64, isSD bool) LxnsMaimaiRequestUserReferBestSong {
	var getReferType string
	if isSD {
		getReferType = "standard"
	} else {
		getReferType = "dx"
	}
	getData, err := web.RequestDataWithHeaders(web.NewDefaultClient(), "https://maimai.lxns.net/api/v0/maimai/player/"+strconv.FormatInt(friendID, 10)+"/bests?song_id="+strconv.FormatInt(songID, 10)+"&song_type="+getReferType, "GET", func(request *http.Request) error {
		request.Header.Add("Authorization", os.Getenv("lxnskey"))
		return nil
	}, nil)
	if err != nil {
		return LxnsMaimaiRequestUserReferBestSong{Success: false}
	}
	var handlerData LxnsMaimaiRequestUserReferBestSong
	json.Unmarshal(getData, &handlerData)
	return handlerData
}

func RequestReferSongIndex(friendID int64, songID int64, diff int64, isSD bool) LxnsMaimaiRequestUserReferBestSongIndex {
	var getReferType string
	if isSD {
		getReferType = "standard"
	} else {
		getReferType = "dx"
	}
	getData, err := web.RequestDataWithHeaders(web.NewDefaultClient(), "https://maimai.lxns.net/api/v0/maimai/player/"+strconv.FormatInt(friendID, 10)+"/bests?song_id="+strconv.FormatInt(songID, 10)+"&song_type="+getReferType+"&level_index="+strconv.FormatInt(diff, 10), "GET", func(request *http.Request) error {
		request.Header.Add("Authorization", os.Getenv("lxnskey"))
		return nil
	}, nil)
	if err != nil {
		return LxnsMaimaiRequestUserReferBestSongIndex{Success: false}
	}
	var handlerData LxnsMaimaiRequestUserReferBestSongIndex
	json.Unmarshal(getData, &handlerData)
	return handlerData
}

func simpleNumHandler(num int, upper bool) int {
	if upper {
		if num < 1000 && num > 100 {
			toint, _ := strconv.Atoi(fmt.Sprintf("10%d", num))
			return toint
		}
		if num > 1000 && num < 10000 {
			toint, _ := strconv.Atoi(fmt.Sprintf("1%d", num))
			return toint
		}
	} else {
		getFmt := fmt.Sprintf("%d", num)
		getFmt = getFmt[2:]
		toint, _ := strconv.Atoi(getFmt)
		return toint
	}
	return num
}

// UpdateHandler Update handler
func UpdateHandler(userMusicList UserMusicListStruct, getTokenId string) int {
	getFullDataStruct := convert(userMusicList)
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
	if err != nil {
		panic(err)
	}
	return resp.StatusCode
}
