package main

import (
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	_ "github.com/wdvxdr1123/emori/pkg/anime"
	_ "github.com/wdvxdr1123/emori/pkg/rcnb"
	_ "github.com/wdvxdr1123/emori/pkg/thisdoesnotexist"
)

func init() {
	log.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[emori][%time%][%lvl%]: %msg% \n",
	})
	log.SetLevel(log.InfoLevel)
}

func main() {
	zero.Run(zero.Config{
		NickName:      []string{""},
		CommandPrefix: ".",
		SuperUsers:    nil,
		Driver: []zero.Driver{
			driver.NewWebSocketClient("127.0.0.1", "6700", ""),
		},
	})
	select {}
}
