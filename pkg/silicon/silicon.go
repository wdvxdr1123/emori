package silicon

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	jsoniter "github.com/json-iterator/go"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension"
	"github.com/wdvxdr1123/ZeroBot/message"
)

type resp struct {
	Code int         `json:"code"`
	Err  interface{} `json:"err"`
	Url  string      `json:"url"`
}

type format struct {
	Language   string `json:"language"`
	Theme      string `json:"theme"`
	LinePad    int    `json:"line_pad"`
	LineOffset int    `json:"line_offset"`
	TabWidth   int    `json:"tab_width"`
}

type req struct {
	Code   string `json:"code"`
	Format format `json:"format"`
}

func init() {
	zero.OnCommand("silicon").Handle(func(ctx *zero.Ctx) {
		m := extension.CommandModel{}
		_ = ctx.Parse(&m)
		args := strings.SplitN(m.Args, "\n", 2)
		for i := range args {
			args[i] = strings.TrimSpace(args[i])
		}
		fmt.Println(args[0])
		buf := &bytes.Buffer{}
		_ = jsoniter.NewEncoder(buf).Encode(&req{
			Code: args[1],
			Format: format{
				Language:   args[0],
				Theme:      "Dracula",
				LinePad:    2,
				LineOffset: 1,
				TabWidth:   4,
			},
		})
		postResp, err := http.Post("https://api.wdvxdr.com/silicon", "application/json", buf)
		if err != nil {
			return
		}
		defer func() { _ = postResp.Body.Close() }()

		var rsp = &resp{}
		_ = jsoniter.NewDecoder(postResp.Body).Decode(rsp)
		if rsp.Code != 200 {
			ctx.SendChain(message.Text(rsp.Err))
		}
		ctx.Send(message.Image(rsp.Url))
	})
}
