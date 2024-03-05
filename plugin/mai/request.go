package mai

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/FloatTech/floatbox/web"
	"io"
	"net/http"
	"os"
)

type DivingFishB50 struct {
	QQ       int    `json:"qq"`
	Username string `json:"username"`
	B50      bool   `json:"b50"`
}

type DivingFishB50UserName struct {
	Username string `json:"username"`
	B50      bool   `json:"b50"`
}

type DivingFishB40 struct {
	QQ       int    `json:"qq"`
	Username string `json:"username"`
}

type DivingFishDevFullDataRecords struct {
	AdditionalRating int    `json:"additional_rating"`
	Nickname         string `json:"nickname"`
	Plate            string `json:"plate"`
	Rating           int    `json:"rating"`
	Records          []struct {
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
	} `json:"records"`
	Username string `json:"username"`
}

func QueryMaiBotDataFromQQ(qq int, isB50 bool) (playerdata []byte, err error) {
	// packed json and sent.
	var jsonStruct interface{}
	if isB50 == true {
		jsonStruct = DivingFishB50{QQ: qq, B50: true}
	} else {
		jsonStruct = DivingFishB40{QQ: qq}
	}
	jsonStructData, err := json.Marshal(jsonStruct)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://www.diving-fish.com/api/maimaidxprober/query/player", bytes.NewBuffer(jsonStructData))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 400 {
		return nil, errors.New("- 未找到用户或者用户数据丢失\n\n - 请检查您是否在 https://www.diving-fish.com/maimaidx/prober/ 上 上传过成绩并且有绑定QQ号 \n\n- 指令为!mai 而不是！mai b50 \n - 上传成绩可通过 https://maimai.bakapiano.com 上传")
	}
	if resp.StatusCode == 403 {
		return nil, errors.New("- 该用户设置禁止查分\n\n - 请检查您是否在 https://www.diving-fish.com/maimaidx/prober/ 上 是否关闭了允许他人查分功能")
	}
	playerGetData, err := io.ReadAll(resp.Body)
	return playerGetData, err
}
func QueryMaiBotDataFromUserName(username string) (playerdata []byte, err error) {
	// packed json and sent.
	jsonStruct := DivingFishB50UserName{Username: username, B50: true}
	jsonStructData, err := json.Marshal(jsonStruct)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://www.diving-fish.com/api/maimaidxprober/query/player", bytes.NewBuffer(jsonStructData))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 400 {
		return nil, errors.New("- 未找到用户或者用户数据丢失\n\n - 请检查您是否在 https://www.diving-fish.com/maimaidx/prober/ 上 上传过成绩并且有绑定QQ号\n- 指令为!mai 而不是！mai b50 \n - 上传成绩可通过 https://maimai.bakapiano.com 上传")
	}
	if resp.StatusCode == 403 {
		return nil, errors.New("- 该用户设置禁止查分\n\n - 请检查您是否在 https://www.diving-fish.com/maimaidx/prober/ 上 是否关闭了允许他人查分功能")
	}
	playerDataByte, err := io.ReadAll(resp.Body)
	return playerDataByte, err
}
func QueryChunDataFromQQ(qq int) (playerdata []byte, err error) {
	// packed json and sent.
	jsonStruct := DivingFishB50{QQ: qq, B50: true}
	jsonStructData, err := json.Marshal(jsonStruct)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "https://www.diving-fish.com/api/chunithmprober/query/player", bytes.NewBuffer(jsonStructData))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		panic(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == 400 {
		return nil, errors.New("- 未找到用户或者用户数据丢失\n\n - 请检查您是否在 https://www.diving-fish.com/maimaidx/prober/ 上 上传过成绩并且有绑定QQ号 \n - 上传成绩可通过 https://maimai.bakapiano.com 上传")
	}
	if resp.StatusCode == 403 {
		return nil, errors.New("- 该用户设置禁止查分\n\n - 请检查您是否在 https://www.diving-fish.com/maimaidx/prober/ 上 是否关闭了允许他人查分功能")
	}
	playerData, err := io.ReadAll(resp.Body)
	return playerData, err
}

func QueryDevDataFromDivingFish(qq string) DivingFishDevFullDataRecords {
	getData, err := web.RequestDataWithHeaders(web.NewDefaultClient(), "https://www.diving-fish.com/api/maimaidxprober/dev/player/records?qq="+qq, "GET", func(request *http.Request) error {
		request.Header.Add("Developer-Token", os.Getenv("dvkey"))
		return nil
	}, nil)
	if err != nil {
		return DivingFishDevFullDataRecords{}
	}
	var handlerData DivingFishDevFullDataRecords
	json.Unmarshal(getData, &handlerData)
	return handlerData
}
