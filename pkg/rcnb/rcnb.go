package rcnb

import (
	"strings"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension"
	"github.com/wdvxdr1123/rcnb.go"
)

var _ = zero.OnCommandGroup([]string{"rcnb", "nbrc"}).SetPriority(10).SetBlock(true).Handle(
	func(ctx *zero.Ctx) {
		m := extension.CommandModel{}
		err := ctx.Parse(&m)
		if err != nil {
			return
		}
		if m.Args == "" {
			ctx.Send("RCNB~")
			return
		}
		switch m.Command {
		case "rcnb":
			v, err := rcnb.Encode([]byte(m.Args))
			if err != nil {
				ctx.Send("由于RC过于NB导致无法编码: " + err.Error())
			}
			ctx.Send("rcnb://" + v)
		case "nbrc":
			args := strings.TrimPrefix(m.Args, "rcnb://")
			v, err := rcnb.Decode(args)
			if err != nil {
				ctx.Send("由于RC过于NB导致无法解码: " + err.Error())
				break
			}
			ctx.Send(string(v))
		}
		return
	},
)
