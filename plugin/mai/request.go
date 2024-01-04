package mai

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
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

func QueryMaiBotDataFromQQ(qq int) (playerdata []byte, err error) {
	// packed json and sent.
	jsonStruct := DivingFishB50{QQ: qq, B50: true}
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
