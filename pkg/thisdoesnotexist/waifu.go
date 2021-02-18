package thisdoesnotexist

import (
	"fmt"
	"math/rand"
	"time"

	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var _ = zero.OnCommandGroup([]string{"thiswaifudoesnotexist", "this_waifu_does_not_exist"}).
	SetBlock(true).
	SetPriority(10).
	Handle(
		func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			const totalImages = 100000
			id := rand.Intn(totalImages)
			zero.Send(event, message.Image(fmt.Sprintf("https://www.thiswaifudoesnotexist.net/example-%v.jpg", id)))
			return zero.FinishResponse
		},
	)

var _ = zero.OnCommandGroup([]string{"thisanimedoesnotexist", "this_anime_does_not_exist"}).
	SetBlock(true).
	SetPriority(10).
	Handle(
		func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			const totalImages = 100000
			id := rand.Intn(totalImages)
			zero.Send(event, message.Image(fmt.Sprintf("https://thisanimedoesnotexist.ai/results/psi-1.0/seed%05d.png", id)))
			return zero.FinishResponse
		},
	)

// https://thisfursonadoesnotexist.com/v2/jpgs-2x/seed61811.jpg
var _ = zero.OnCommandGroup([]string{"thisfursonadoesnotexist", "this_fursona_does_not_exist"}).
	SetBlock(true).
	SetPriority(10).
	Handle(
		func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			const totalImages = 100000
			id := rand.Intn(totalImages)
			zero.Send(event, message.Image(fmt.Sprintf("https://thisfursonadoesnotexist.com/v2/jpgs-2x/seed%05d.jpg", id)))
			return zero.FinishResponse
		},
	)

// https://thisponydoesnotexist.net/v1/w2x-redo/jpgs/seed29775.jpg
var _ = zero.OnCommandGroup([]string{"thisponydoesnotexist", "this_pony_does_not_exist"}).
	SetBlock(true).
	SetPriority(10).
	Handle(
		func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			const totalImages = 100000
			id := rand.Intn(totalImages)
			zero.Send(event, message.Image(fmt.Sprintf("https://thisponydoesnotexist.net/v1/w2x-redo/jpgs/seed%05d.jpg", id)))
			return zero.FinishResponse
		},
	)

// https://thiscatdoesnotexist.com/
var _ = zero.OnCommandGroup([]string{"thiscatdoesnotexist", "this_cat_does_not_exist"}).
	SetBlock(true).
	SetPriority(10).
	Handle(
		func(matcher *zero.Matcher, event zero.Event, state zero.State) zero.Response {
			zero.Send(event, message.MessageSegment{
				Type: "image",
				Data: map[string]string{
					"file":  "https://thiscatdoesnotexist.com",
					"cache": "0",
				},
			})
			return zero.FinishResponse
		},
	)
