package mai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/gg"
	"github.com/FloatTech/imgfactory"
	zero "github.com/wdvxdr1123/ZeroBot"
	"golang.org/x/text/width"
)

type LxnsMaimaiRequestFromQQ struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Name       string `json:"name"`
		Rating     int    `json:"rating"`
		FriendCode int64  `json:"friend_code"`
		Trophy     struct {
			Id    int    `json:"id"`
			Name  string `json:"name"`
			Color string `json:"color"`
		} `json:"trophy"`
		CourseRank int    `json:"course_rank"`
		ClassRank  int    `json:"class_rank"`
		Star       int    `json:"star"`
		IconUrl    string `json:"icon_url"`
		NamePlate  struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"name_plate"`
		Frame struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"frame"`
		UploadTime time.Time `json:"upload_time"`
	} `json:"data"`
}

type LxnsMaimaiRequestB50 struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		StandardTotal int `json:"standard_total"`
		DxTotal       int `json:"dx_total"`
		Standard      []struct {
			Id           int       `json:"id"`
			SongName     string    `json:"song_name"`
			Level        string    `json:"level"`
			LevelIndex   int       `json:"level_index"`
			Achievements float64   `json:"achievements"`
			Fc           *string   `json:"fc"`
			Fs           *string   `json:"fs"`
			DxScore      int       `json:"dx_score"`
			DxRating     float64   `json:"dx_rating"`
			Rate         string    `json:"rate"`
			Type         string    `json:"type"`
			UploadTime   time.Time `json:"upload_time"`
		} `json:"standard"`
		Dx []struct {
			Id           int       `json:"id"`
			SongName     string    `json:"song_name"`
			Level        string    `json:"level"`
			LevelIndex   int       `json:"level_index"`
			Achievements float64   `json:"achievements"`
			Fc           *string   `json:"fc"`
			Fs           *string   `json:"fs"`
			DxScore      int       `json:"dx_score"`
			DxRating     float64   `json:"dx_rating"`
			Rate         string    `json:"rate"`
			Type         string    `json:"type"`
			UploadTime   time.Time `json:"upload_time"`
		} `json:"dx"`
	} `json:"data"`
}

type LxnsMaimaiRequestDataPiece struct {
	Id           int       `json:"id"`
	SongName     string    `json:"song_name"`
	Level        string    `json:"level"`
	LevelIndex   int       `json:"level_index"`
	Achievements float64   `json:"achievements"`
	Fc           *string   `json:"fc"`
	Fs           *string   `json:"fs"`
	DxScore      int       `json:"dx_score"`
	DxRating     float64   `json:"dx_rating"`
	Rate         string    `json:"rate"`
	Type         string    `json:"type"`
	UploadTime   time.Time `json:"upload_time"`
}

// NOTE: lxns network maimai uses qq => friendCode

// RequestBasicDataFromLxns

/*
	{
	    "success": true,
	    "code": 200,
	    "data": {
	        "name": "StarKooi",
	        "rating": 11616,
	        "friend_code": 00000000000,
	        "trophy": {
	            "id": 258509,
	            "name": "きみもヴァンパイア",
	            "color": "Normal"
	        },
	        "course_rank": 5,
	        "class_rank": 1,
	        "star": 64,
	        "icon_url": "https://maimai.wahlap.com/maimai-mobile/img/Icon/c3289d7ae91077ac.png",
	        "name_plate": {
	            "id": 255901,
	            "name": "すりぃちほー"
	        },
	        "frame": {
	            "id": 250701,
	            "name": "ホワイトボード(はっぴー)"
	        },
	        "upload_time": "2024-01-02T08:42:48Z"
	    }
	}
*/

func RequestBasicDataFromLxns(qq int64) LxnsMaimaiRequestFromQQ {
	getData, err := web.RequestDataWithHeaders(web.NewDefaultClient(), "https://maimai.lxns.net/api/v0/maimai/player/qq/"+strconv.FormatInt(qq, 10), "GET", func(request *http.Request) error {
		request.Header.Add("Authorization", os.Getenv("lxnskey"))
		return nil
	}, nil)
	if err != nil {
		return LxnsMaimaiRequestFromQQ{}
	}
	var handlerData LxnsMaimaiRequestFromQQ
	json.Unmarshal(getData, &handlerData)
	return handlerData
}

func RequestB50DataByFriendCode(friendCode int64) LxnsMaimaiRequestB50 {
	getData, err := web.RequestDataWithHeaders(web.NewDefaultClient(), "https://maimai.lxns.net/api/v0/maimai/player/"+strconv.FormatInt(friendCode, 10)+"/bests", "GET", func(request *http.Request) error {
		request.Header.Add("Authorization", os.Getenv("lxnskey"))
		return nil
	}, nil)
	if err != nil {
		return LxnsMaimaiRequestB50{}
	}
	var handlerData LxnsMaimaiRequestB50
	json.Unmarshal(getData, &handlerData)

	return handlerData
}

func ReCardRenderBase(data LxnsMaimaiRequestDataPiece, getNum int) image.Image {
	getType := data.Type
	var CardBackGround string
	var multiTypeRender sync.WaitGroup
	var CoverDownloader sync.WaitGroup
	CoverDownloader.Add(1)
	multiTypeRender.Add(1)
	// choose Type.
	if getType == "standard" {
		CardBackGround = typeImageSD
	} else {
		CardBackGround = typeImageDX
	}
	charCount := 0.0
	setBreaker := false
	var truncated string
	var charFloatNum float64
	getSongName := data.SongName
	getSongId := strconv.Itoa(data.Id)
	var Image image.Image
	go func() {
		defer CoverDownloader.Done()
		Image, _ = GetCoverLxns(getSongId)
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
	drawBackGround := gg.NewContextForImage(ReturnMaiIndexBackground(data.LevelIndex))
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
	drawBackGround.DrawString("#"+strconv.Itoa(getNum), 130, 111)
	drawBackGround.FillPreserve()
	// draw rest of card.
	drawBackGround.SetFontFace(levelFont)
	getCount := GetShouldCount(data.Achievements)
	actuallyOutput := data.DxRating / getCount * 100 / data.Achievements
	actuallyOutputF := strconv.FormatFloat(actuallyOutput, 'f', 1, 64)
	drawBackGround.DrawString(actuallyOutputF, 195, 111)
	drawBackGround.FillPreserve()
	drawBackGround.SetFontFace(ratingFont)
	drawBackGround.DrawString("▶", 235, 111)
	drawBackGround.FillPreserve()
	drawBackGround.SetFontFace(ratingFont)
	drawBackGround.DrawString(strconv.Itoa(data.DxScore), 250, 111)
	drawBackGround.FillPreserve()
	if data.Fc != nil {
		FcPointer := *data.Fc
		drawBackGround.DrawImage(LoadComboImage(FcPointer), 290, 84)
	}
	if data.Fs != nil {
		FsPointer := *data.Fs
		drawBackGround.DrawImage(LoadSyncImage(FsPointer), 325, 84)
	}
	drawBackGround.DrawImage(loadSongType, 68, 88)
	return drawBackGround.Image()
}

func DataPiecesRepacked(data LxnsMaimaiRequestB50, returnTypeIsSD bool, getShouldNum int) LxnsMaimaiRequestDataPiece {
	if returnTypeIsSD {
		return LxnsMaimaiRequestDataPiece{
			Id:           data.Data.Standard[getShouldNum].Id,
			SongName:     data.Data.Standard[getShouldNum].SongName,
			Level:        data.Data.Standard[getShouldNum].Level,
			LevelIndex:   data.Data.Standard[getShouldNum].LevelIndex,
			Achievements: data.Data.Standard[getShouldNum].Achievements,
			Fc:           data.Data.Standard[getShouldNum].Fc,
			Fs:           data.Data.Standard[getShouldNum].Fs,
			DxScore:      data.Data.Standard[getShouldNum].DxScore,
			DxRating:     data.Data.Standard[getShouldNum].DxRating,
			Rate:         data.Data.Standard[getShouldNum].Rate,
			Type:         data.Data.Standard[getShouldNum].Type,
			UploadTime:   data.Data.Standard[getShouldNum].UploadTime,
		}
	} else {
		return LxnsMaimaiRequestDataPiece{
			Id:           data.Data.Dx[getShouldNum].Id,
			SongName:     data.Data.Dx[getShouldNum].SongName,
			Level:        data.Data.Dx[getShouldNum].Level,
			LevelIndex:   data.Data.Dx[getShouldNum].LevelIndex,
			Achievements: data.Data.Dx[getShouldNum].Achievements,
			Fc:           data.Data.Dx[getShouldNum].Fc,
			Fs:           data.Data.Dx[getShouldNum].Fs,
			DxScore:      data.Data.Dx[getShouldNum].DxScore,
			DxRating:     data.Data.Dx[getShouldNum].DxRating,
			Rate:         data.Data.Dx[getShouldNum].Rate,
			Type:         data.Data.Dx[getShouldNum].Type,
			UploadTime:   data.Data.Dx[getShouldNum].UploadTime,
		}
	}

}

func ReFullPageRender(data LxnsMaimaiRequestB50, userData LxnsMaimaiRequestFromQQ, ctx *zero.Ctx) (image.Image, bool) {
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
			// check again if user use origin plate
			if userData.Data.NamePlate.Id != 0 {
				getBackground, err := web.GetData("https://lxns.org/maimai/plate/" + strconv.FormatInt(int64(userData.Data.NamePlate.Id), 10) + ".png")
				if err != nil {
					b50bg = b50bgOriginal
				}
				// resize pic
				b50bg = b50Custom
				getImage, _, _ := image.Decode(bytes.NewReader(getBackground))
				images := Resize(getImage, 1260, 210)
				b50Render.DrawImage(images, 595, 30)
				b50Render.Fill()
			} else {
				b50bg = b50bgOriginal
			}
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
	b50Render.DrawStringAnchored(width.Widen.String(userData.Data.Name), 825, 160, 0, 0)
	b50Render.Fill()
	b50Render.SetFontFace(titleFont)
	setPlateLocalStatus := GetUserInfoFromDatabase(ctx.Event.UserID)
	var dataPlate bool
	// tips trophy was custom plate here
	if setPlateLocalStatus != "" {
		userData.Data.Trophy.Name = setPlateLocalStatus
		dataPlate = true
	} else {
		dataPlate = false
	}
	b50Render.DrawStringAnchored(strings.Join(strings.Split(userData.Data.Trophy.Name, ""), " "), 1050, 207, 0.5, 0.5)
	b50Render.Fill()
	getRating := getRatingBg(userData.Data.Rating)
	getRatingBG, err := gg.LoadImage(loadMaiPic + getRating)
	if err != nil {
		panic(err)
	}
	b50Render.DrawImage(getRatingBG, 800, 40)
	b50Render.Fill()
	// render Rank
	imgs, err := GetRankPicRaw(userData.Data.CourseRank)
	if err != nil {
		panic(err)
	}
	b50Render.DrawImage(imgs, 1080, 50)
	b50Render.Fill()
	// draw number
	b50Render.SetFontFace(scoreFont)
	b50Render.SetRGBA255(236, 219, 113, 255)
	b50Render.DrawStringAnchored(strconv.Itoa(userData.Data.Rating), 1056, 60, 1, 1)
	b50Render.Fill()
	// Render Card Type
	getSDLength := len(data.Data.Standard)
	fmt.Printf(strconv.Itoa(getSDLength) + "\n")
	getDXLength := len(data.Data.Dx)
	getDXinitX := 45
	getDXinitY := 1225
	getInitX := 45
	getInitY := 285
	var i int
	for i = 0; i < getSDLength; i++ {
		b50Render.DrawImage(ReCardRenderBase(DataPiecesRepacked(data, true, i), i+1), getInitX, getInitY)
		getInitX += 400
		if getInitX == 2045 {
			getInitX = 45
			getInitY += 125
		}
	}

	for dx := 0; dx < getDXLength; dx++ {
		b50Render.DrawImage(ReCardRenderBase(DataPiecesRepacked(data, false, dx), dx+1), getDXinitX, getDXinitY)
		getDXinitX += 400
		if getDXinitX == 2045 {
			getDXinitX = 45
			getDXinitY += 125
		}
	}
	return b50Render.Image(), dataPlate
}

// GetCoverLxns Careful The nil data
func GetCoverLxns(id string) (image.Image, error) {
	fileName := id + ".png"
	filePath := Root + "coverLxns/" + fileName
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Auto download cover from diving fish's site
		downloadURL := "https://lxns.org/maimai/jacket/" + fileName
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

func ReturnMaiIndexBackground(returnInt int) image.Image {
	var chart string
	switch {
	case returnInt == 0:
		chart = "basic"
	case returnInt == 1:
		chart = "advanced"
	case returnInt == 2:
		chart = "expert"
	case returnInt == 3:
		chart = "master"
	default:
		chart = "remaster"
	}
	data, _ := gg.LoadImage(loadMaiPic + "chart_" + NoHeadLineCase(chart) + ".png")
	return data
}

func GetShouldCount(archivement float64) float64 {
	switch {
	case archivement >= 100.5:
		return 22.4
	case archivement >= 100.0:
		return 21.6
	case archivement >= 99.5:
		return 21.1
	case archivement >= 99.0:
		return 20.8
	case archivement >= 98.0:
		return 20.3
	case archivement >= 97.0:
		return 20.0
	case archivement >= 94.0:
		return 16.8
	case archivement >= 90.0:
		return 13.6
	case archivement >= 80.0:
		return 12.8
	case archivement >= 75.0:
		return 12.0
	case archivement >= 70.0:
		return 11.2
	case archivement >= 60.0:
		return 9.6
	case archivement >= 50.0:
		return 8.0
	case archivement >= 40.0:
		return 6.4
	case archivement >= 30.0:
		return 4.8
	case archivement >= 20.0:
		return 3.2
	case archivement >= 10.0:
		return 1.6
	default:
		return 0.0
	}
}
