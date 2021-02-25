package anime

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func ReplyRule(messageID int64) zero.Rule {
	var mid = strconv.FormatInt(messageID, 10)
	return func(ctx *zero.Ctx) bool {
		if len(ctx.Event.Message) <= 0 {
			return false
		}
		if ctx.Event.Message[0].Type != "reply" {
			return false
		}
		return ctx.Event.Message[0].Data["id"] == mid
	}
}

var _ = zero.OnCommandGroup([]string{"searchtag", "search_tag", "添加标签"}).Handle(
	func(ctx *zero.Ctx) {
		var model = &extension.CommandModel{}
		err := ctx.Parse(model)
		if err != nil {
			ctx.Send(err.Error())
			return
		}
		tags, err := GetTagByKeyword(model.Args)
		var msg = "为您找到以下标签:\n"
		for i, v := range tags {
			msg += fmt.Sprintf("%v: %v\n", i, v.Name)
		}
		msg += `如需添加标签，请回复该消息并附带编号，2分钟内有效!`
		sid := ctx.Send(strings.TrimSpace(msg))
		recv, cancel := ctx.FutureEvent("message", ReplyRule(sid)).Repeat()
		timeout := time.After(2 * time.Minute)
		for {
			select {
			case e := <-recv:
				arg := strings.TrimSpace(e.Message.ExtractPlainText())
				if arg == "ok" {
					cancel()
					return
				}
				i, err := strconv.ParseInt(arg, 10, 64)
				if err != nil || i < 0 || int(i) >= len(tags) {
					ctx.Send("参数无效,请重新输入!")
					continue
				}
				addTags(e.UserID, tags[i])
				ctx.Send(message.Message{
					message.Reply(e.MessageID),
					message.Text(fmt.Sprint("已为您添加Tag: ", tags[i].Name, " !")),
				})
			case <-timeout:
				cancel()
				return
			}
		}
	},
)

var _ = zero.OnCommandGroup([]string{"fetch_annie", "fetch"}, zero.OnlyGroup).Handle(
	func(ctx *zero.Ctx) {
		var cm = extension.CommandModel{}
		err := ctx.Parse(&cm)
		yi, _ := strconv.Atoi(cm.Args)
		if yi <= 0 {
			yi = 1
		}
		if err != nil {
			ctx.Send(fmt.Sprint("消息处理失败: ", err))
			return
		}
		tags := queryTags(ctx.Event.UserID)
		annie, err := getAnnieInfo(tags)
		siz := len(annie)
		totY := siz / 5
		if siz%5 != 0 {
			totY++
		}
		if err != nil {
			ctx.Send(fmt.Sprint("无法获取资源信息", err))
			return
		}
		var msg = message.Message{
			message.CustomNode("Anime", ctx.Event.SelfID, fmt.Sprintf(`总共记录%v条
当前页码：%v
总页码: %v
为您找到以下资源:`, siz, yi, totY)),
		}
		for i := (yi - 1) * 5; i < yi*5 && i < siz; i++ {
			msg = append(msg, message.CustomNode(
				"Anime",
				ctx.Event.SelfID,
				fmt.Sprintf("编号: %v\n标题：%v\n链接: %v", i, annie[i].title, annie[i].link),
			))
		}
		msg = append(msg, message.CustomNode(
			"Anime",
			ctx.Event.SelfID,
			`如需观看视频，请回复该消息并附带编号，2分钟内有效!`,
		))
		sid := ctx.SendGroupForwardMessage(ctx.Event.GroupID, msg).Get("message_id").Int()
		recv, cancel := ctx.FutureEvent("message", ReplyRule(sid)).Repeat()
		timeout := time.After(2 * time.Minute)
		for {
			select {
			case e := <-recv:
				arg := strings.TrimSpace(e.Message.ExtractPlainText())
				if arg == "ok" {
					cancel()
					return
				}
				i, err := strconv.ParseInt(arg, 10, 64)
				if err != nil || i < 0 || int(i) >= len(annie) {
					ctx.Send("参数无效,请重新输入!")
					continue
				}
				ctx.Send("正在处理，请稍等几分钟...")
				file, err := DownloadAnnie(annie[i].torrentLink)
				if err != nil {
					ctx.Send(fmt.Sprint("下载失败: ", err))
					continue
				}
				fileList, err := spilitVideo(file, 100)
				if err != nil {
					ctx.Send(fmt.Sprint("切分视频失败: ", err))
					continue
				}
				var msg = message.Message{}
				for _, v := range fileList {
					msg = append(msg, message.CustomNode(
						"Anime",
						ctx.Event.SelfID,
						fmt.Sprintf("[CQ:video,file=file:///%v]", v),
					))
				}
				ctx.SendGroupForwardMessage(ctx.Event.GroupID, msg)
				cancel()
				return
			case <-timeout:
				cancel()
				return
			}
		}
	},
)

var _ = zero.OnCommandGroup([]string{"mytag", "my_tag"}).Handle(
	func(ctx *zero.Ctx) {
		tags := queryTags(ctx.Event.UserID)
		var msg = "您当前已添加以下标签:\n"
		for i, v := range tags {
			msg += fmt.Sprintf("%v: %v\n", i, v.Name)
		}
		msg += `如需删除标签，请回复该消息并附带编号，2分钟内有效!`
		sid := ctx.Send(strings.TrimSpace(msg))
		recv, cancel := ctx.FutureEvent("message", ReplyRule(sid)).Repeat()
		timeout := time.After(2 * time.Minute)
		for {
			select {
			case e := <-recv:
				arg := strings.TrimSpace(e.Message.ExtractPlainText())
				if arg == "ok" {
					cancel()
					return
				}
				i, err := strconv.ParseInt(arg, 10, 64)
				if err != nil || i < 0 || int(i) >= len(tags) {
					ctx.Send("参数无效,请重新输入!")
					continue
				}
				deleteTags(e.UserID, tags[i])
				ctx.Send(message.Message{
					message.Reply(e.MessageID),
					message.Text(fmt.Sprint("删除成功!")),
				})
			case <-timeout:
				cancel()
				return
			}
		}
	},
)
