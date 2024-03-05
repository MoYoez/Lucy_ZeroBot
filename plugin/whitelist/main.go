package whitelist

import (
	"encoding/json"
	"github.com/FloatTech/floatbox/file"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"log"
	"os"
	"strconv"
)

type whitelist struct {
	Or []struct {
		GroupId struct {
			In []int `json:".in"`
		} `json:"group_id,omitempty"`
		UserId interface{} `json:"user_id"`
	} `json:".or"`
}

var (
	engine = control.Register("whitelist", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Help:             "Hi NekoPachi!\n",
	})
)

func init() {
	engine.OnRegex(`^/whitelist\s(.*)`, zero.SuperUserPermission).SetBlock(true).Handle(func(ctx *zero.Ctx) {
		getNumber, err := strconv.ParseInt(ctx.State["regex_matched"].([]string)[1], 10, 64)
		if err != nil {
			panic(err)
		}
		reader, err := os.ReadFile(file.BOTPATH + "/filter.json")
		if err != nil {
			panic(err)
		}
		var list whitelist
		_ = json.Unmarshal(reader, &list)
		list.Or[0].GroupId.In = append(list.Or[0].GroupId.In, int(getNumber))
		modifiedData, err := json.MarshalIndent(list, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		err = os.WriteFile(file.BOTPATH+"/filter.json", modifiedData, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		ctx.CallAction("reload_event_filter", zero.Params{"file": "filter.json"})
		ctx.SendChain(message.Text("执行完毕~"))
	})

}
