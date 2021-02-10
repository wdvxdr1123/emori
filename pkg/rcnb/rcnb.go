package rcnb

import (
	"strconv"
	"strings"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension"
	"github.com/wdvxdr1123/ZeroBot/message"
	"github.com/wdvxdr1123/rcnb.go"
)

var _ = zero.OnCommandGroup([]string{"rcnb", "nbrc"}).SetPriority(10).SetBlock(true).Handle(
	func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
		m := extension.CommandModel{}
		err := state.Parse(&m)
		if err != nil {
			return zero.FinishResponse
		}
		if m.Args == "" {
			zero.Send(event, "RCNB~")
			return zero.FinishResponse
		}
		switch m.Command {
		case "rcnb":
			v, err := rcnb.Encode([]byte(m.Args))
			if err != nil {
				zero.Send(event, "由于RC过于NB导致无法编码: "+err.Error())
			}
			zero.Send(event, message.Message{
				message.Reply(strconv.FormatInt(event.MessageID, 10)),
				message.Text("rcnb://" + v),
			})
		case "nbrc":
			args := strings.TrimPrefix(m.Args, "rcnb://")
			v, err := rcnb.Decode(args)
			if err != nil {
				zero.Send(event, "由于RC过于NB导致无法解码: "+err.Error())
			}
			zero.Send(event, message.Message{
				message.Reply(strconv.FormatInt(event.MessageID, 10)),
				message.Text(string(v)),
			})
		}
		return zero.FinishResponse
	},
)
