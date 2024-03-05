// Package funwork for fun
package funwork

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/FloatTech/floatbox/math"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

func init() {
	engine.OnFullMatchGroup([]string{"今日本群RBQ"}, zero.OnlyGroup).SetBlock(false).Limit(ctxext.LimitByGroup).
		Handle(func(ctx *zero.Ctx) {
			list := ctx.CallAction("get_group_member_list", zero.Params{
				"group_id": ctx.Event.GroupID,
				"no_cache": false,
			}).Data
			temp := list.Array()
			sort.SliceStable(temp, func(i, j int) bool {
				return temp[i].Get("last_sent_time").Int() < temp[j].Get("last_sent_time").Int()
			})
			temp = temp[math.Max(0, len(temp)-30):]
			now := time.Now()
			s := md5.Sum(helper.StringToBytes(fmt.Sprintf("%d%d%d", now.Year(), now.Month(), now.Day())))
			r := rand.New(rand.NewSource(int64(binary.LittleEndian.Uint64(s[:]))))
			who := temp[r.Intn(len(temp))]
			if who.Get("user_id").Int() == 2854196310 {
				// re rand.
				who = temp[r.Intn(len(temp))]
				if who.Get("user_id").Int() == 2854196310 {
					ctx.SendChain(message.Text("今天没有群RBQ惹 aww"))
					return
				}
			}
			cutename := who.Get("card").Str
			cuteid := who.Get("user_id").Int()
			if cutename == "" {
				cutename = who.Get("nickname").Str
			}
			avtar := fmt.Sprintf("[CQ:image,file=https://q4.qlogo.cn/g?b=qq&nk=%d&s=640,cache=0]", cuteid)
			msg := fmt.Sprintf("今日本群RBQ是%s\n[%s](%d)哒！", avtar, cutename, cuteid)
			msg = message.UnescapeCQCodeText(msg)
			ctx.SendGroupMessage(ctx.Event.GroupID, message.ParseMessageFromString(msg))
		})
}
