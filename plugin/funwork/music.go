// Package funwork music 网易云 点歌 Powered By LemonKoi Vercel Services.
package funwork

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/FloatTech/floatbox/web"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/rate"
	"github.com/wdvxdr1123/ZeroBot/message"
)

var limitForMusic = rate.NewManager[int64](time.Minute*3, 8)

type SongJsoner struct {
	Result struct {
		Songs []struct {
			Id      int    `json:"id"`
			Name    string `json:"name"`
			Artists []struct {
				Id        int           `json:"id"`
				Name      string        `json:"name"`
				PicUrl    interface{}   `json:"picUrl"`
				Alias     []interface{} `json:"alias"`
				AlbumSize int           `json:"albumSize"`
				PicId     int           `json:"picId"`
				FansGroup interface{}   `json:"fansGroup"`
				Img1V1Url string        `json:"img1v1Url"`
				Img1V1    int           `json:"img1v1"`
				Trans     interface{}   `json:"trans"`
			} `json:"artists"`
			Album struct {
				Id     int    `json:"id"`
				Name   string `json:"name"`
				Artist struct {
					Id        int           `json:"id"`
					Name      string        `json:"name"`
					PicUrl    interface{}   `json:"picUrl"`
					Alias     []interface{} `json:"alias"`
					AlbumSize int           `json:"albumSize"`
					PicId     int           `json:"picId"`
					FansGroup interface{}   `json:"fansGroup"`
					Img1V1Url string        `json:"img1v1Url"`
					Img1V1    int           `json:"img1v1"`
					Trans     interface{}   `json:"trans"`
				} `json:"artist"`
				PublishTime int64 `json:"publishTime"`
				Size        int   `json:"size"`
				CopyrightId int   `json:"copyrightId"`
				Status      int   `json:"status"`
				PicId       int64 `json:"picId"`
				Mark        int   `json:"mark"`
			} `json:"album"`
			Duration    int           `json:"duration"`
			CopyrightId int           `json:"copyrightId"`
			Status      int           `json:"status"`
			Alias       []interface{} `json:"alias"`
			Rtype       int           `json:"rtype"`
			Ftype       int           `json:"ftype"`
			Mvid        int           `json:"mvid"`
			Fee         int           `json:"fee"`
			RUrl        interface{}   `json:"rUrl"`
			Mark        int           `json:"mark"`
		} `json:"songs"`
		HasMore   bool `json:"hasMore"`
		SongCount int  `json:"songCount"`
	} `json:"result"`
	Code int `json:"code"`
}

func init() {
	engine.OnRegex(`^点歌\s?(.{1,25})$`).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			if !limitForMusic.Load(ctx.Event.GroupID).Acquire() {
				ctx.SendChain(message.Text("太快了哦，麻烦慢一点~"))
				return
			}
			getTimeStamp := strconv.FormatInt(time.Now().Unix(), 10)
			requestURL := "https://nem.lemonkoi.one/search?limit=1&realIP=116.25.146.177&timestamp=" + getTimeStamp + "&keywords=" + url.QueryEscape(ctx.State["regex_matched"].([]string)[1])
			data, err := web.GetData(requestURL)
			var webStatusCode *http.Response
			if err != nil {
				fmt.Print("ERR:", webStatusCode.StatusCode)
			}
			var dataSaver SongJsoner
			_ = json.Unmarshal(data, &dataSaver)
			ctx.SendChain(message.Music("163", int64(dataSaver.Result.Songs[0].Id)))
		})
}
