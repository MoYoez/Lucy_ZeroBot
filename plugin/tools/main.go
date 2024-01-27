// Package tools for tools
package tools

import (
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"

	zero "github.com/wdvxdr1123/ZeroBot"
)

type WhitelistStruct struct {
	Or []struct {
		GroupId struct {
			In []int `json:".in"`
		} `json:"group_id"`
		UserId interface{} `json:"user_id"`
	} `json:".or"`
}

var (
	engine = control.Register("tools", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Help:              "Hi NekoPachi!",
		PrivateDataFolder: "tools",
	})
)

func init() {

}
