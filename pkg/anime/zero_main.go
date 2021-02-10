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
	return func(event *zero.Event, state zero.State) bool {
		if len(event.Message) <= 0 {
			return false
		}
		if event.Message[0].Type != "reply" {
			return false
		}
		return event.Message[0].Data["id"] == mid
	}
}

var _ = zero.OnCommandGroup([]string{"searchtag", "search_tag", "添加标签"}).Handle(
	func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
		var model = &extension.CommandModel{}
		err := state.Parse(model)
		if err != nil {
			zero.Send(event, err)
			return zero.FinishResponse
		}
		tags, err := GetTagByKeyword(model.Args)
		var msg = "为您找到以下标签:\n"
		for i, v := range tags {
			msg += fmt.Sprintf("%v: %v\n", i, v.Name)
		}
		msg += `如需添加标签，请回复该消息并附带编号，2分钟内有效!`
		sid := zero.Send(event, strings.TrimSpace(msg))
		recv, cancel := matcher.FutureEvent("message", ReplyRule(sid)).Repeat()
		timeout := time.After(2 * time.Minute)
		for {
			select {
			case e := <-recv:
				arg := strings.TrimSpace(e.Message.ExtractPlainText())
				if arg == "ok" {
					cancel()
					return zero.FinishResponse
				}
				i, err := strconv.ParseInt(arg, 10, 64)
				if err != nil || i < 0 || int(i) >= len(tags) {
					zero.Send(event, "参数无效,请重新输入!")
					continue
				}
				addTags(e.UserID, tags[i])
				zero.Send(event, message.Message{
					message.Reply(strconv.FormatInt(e.MessageID, 10)),
					message.Text(fmt.Sprint("已为您添加Tag: ", tags[i].Name, " !")),
				})
			case <-timeout:
				cancel()
				return zero.FinishResponse
			}
		}
	},
)

var _ = zero.OnCommandGroup([]string{"fetch_annie", "fetch"}, zero.OnlyGroup).Handle(
	func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
		var cm = extension.CommandModel{}
		err := state.Parse(&cm)
		yi, _ := strconv.Atoi(cm.Args)
		if yi <= 0 {
			yi = 1
		}
		if err != nil {
			zero.Send(event, fmt.Sprint("消息处理失败: ", err))
			return zero.FinishResponse
		}
		tags := queryTags(event.UserID)
		annie, err := getAnnieInfo(tags)
		siz := len(annie)
		totY := siz / 5
		if siz%5 != 0 {
			totY++
		}
		if err != nil {
			zero.Send(event, fmt.Sprint("无法获取资源信息", err))
			return zero.FinishResponse
		}
		var msg = message.Message{
			message.CustomNode("Anime", zero.BotConfig.SelfID, fmt.Sprintf(`总共记录%v条
当前页码：%v
总页码: %v
为您找到以下资源:`, siz, yi, totY)),
		}
		for i := (yi - 1) * 5; i < yi*5 && i < siz; i++ {
			msg = append(msg, message.CustomNode(
				"Anime",
				zero.BotConfig.SelfID,
				fmt.Sprintf("编号: %v\n标题：%v\n链接: %v", i, annie[i].title, annie[i].link),
			))
		}
		msg = append(msg, message.CustomNode(
			"Anime",
			zero.BotConfig.SelfID,
			`如需观看视频，请回复该消息并附带编号，2分钟内有效!`,
		))
		sid := zero.SendGroupForwardMessage(event.GroupID, msg).Get("message_id").Int()
		recv, cancel := matcher.FutureEvent("message", ReplyRule(sid)).Repeat()
		timeout := time.After(2 * time.Minute)
		for {
			select {
			case e := <-recv:
				arg := strings.TrimSpace(e.Message.ExtractPlainText())
				if arg == "ok" {
					cancel()
					return zero.FinishResponse
				}
				i, err := strconv.ParseInt(arg, 10, 64)
				if err != nil || i < 0 || int(i) >= len(annie) {
					zero.Send(event, "参数无效,请重新输入!")
					continue
				}
				zero.Send(event, "正在处理，请稍等几分钟...")
				file, err := DownloadAnnie(annie[i].torrentLink)
				if err != nil {
					zero.Send(event, fmt.Sprint("下载失败: ", err))
					continue
				}
				fileList, err := spilitVideo(file, 100)
				if err != nil {
					zero.Send(event, fmt.Sprint("切分视频失败: ", err))
					continue
				}
				var msg = message.Message{}
				for _, v := range fileList {
					msg = append(msg, message.CustomNode(
						"Anime",
						zero.BotConfig.SelfID,
						fmt.Sprintf("[CQ:video,file=file:///%v]", v),
					))
				}
				zero.SendGroupForwardMessage(event.GroupID, msg)
				cancel()
				return zero.FinishResponse
			case <-timeout:
				cancel()
				return zero.FinishResponse
			}
		}
	},
)

var _ = zero.OnCommandGroup([]string{"mytag", "my_tag"}).Handle(
	func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
		tags := queryTags(event.UserID)
		var msg = "您当前已添加以下标签:\n"
		for i, v := range tags {
			msg += fmt.Sprintf("%v: %v\n", i, v.Name)
		}
		msg += `如需删除标签，请回复该消息并附带编号，2分钟内有效!`
		sid := zero.Send(event, strings.TrimSpace(msg))
		recv, cancel := matcher.FutureEvent("message", ReplyRule(sid)).Repeat()
		timeout := time.After(2 * time.Minute)
		for {
			select {
			case e := <-recv:
				arg := strings.TrimSpace(e.Message.ExtractPlainText())
				if arg == "ok" {
					cancel()
					return zero.FinishResponse
				}
				i, err := strconv.ParseInt(arg, 10, 64)
				if err != nil || i < 0 || int(i) >= len(tags) {
					zero.Send(event, "参数无效,请重新输入!")
					continue
				}
				deleteTags(e.UserID, tags[i])
				zero.Send(event, message.Message{
					message.Reply(strconv.FormatInt(e.MessageID, 10)),
					message.Text(fmt.Sprint("删除成功!")),
				})
			case <-timeout:
				cancel()
				return zero.FinishResponse
			}
		}
	},
)
