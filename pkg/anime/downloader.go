package anime

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sync"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/levigross/grequests"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
)

func DownloadAnnie(t string) (string, error) {
	client, err := torrent.NewClient(torrent.NewDefaultClientConfig())
	if err != nil {
		return "", errors.Wrap(err, "init client failed")
	}
	defer client.Close()
	resp, err := grequests.Get(t, nil)
	if err != nil {
		return "", errors.Wrap(err, "download torrent file failed")
	}
	mi, err := metainfo.Load(resp)
	if err != nil {
		return "", errors.Wrap(err, "parse torrent failed")
	}
	torr, err := client.AddTorrent(mi)
	if err != nil {
		return "", errors.Wrap(err, "add torrent to client failed")
	}
	torr.DownloadAll()
	_, _ = io.Copy(ioutil.Discard, torr.Files()[0].Torrent().NewReader())
	return torr.Files()[0].DisplayPath(), nil
}

func getVideoInfo(fileName string) (uint64, float64, error) {
	// ffprobe -v quiet -print_format json -show_format -i <file>
	cmd := exec.Command(
		"ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", "-i", fileName,
	)
	var b = &bytes.Buffer{}
	cmd.Stdout = b
	err := cmd.Run()
	if err != nil {
		return 0, 0.0, errors.Wrap(err, "run ffprobe failed")
	}
	result := gjson.ParseBytes(b.Bytes())
	return result.Get("format.bit_rate").Uint(), result.Get("format.duration").Float(), nil
}

func spilitVideo(fileName string, maxSize int) ([]string, error) {
	rate, duration, err := getVideoInfo(fileName)
	if err != nil {
		return nil, errors.Wrap(err, "get bit_rate failed")
	}
	fmt.Println(duration)
	// 分片时间 = (文件限制大小[M] * 1024 * 1024) / (媒体文件比特率[b] * 突发码率比率 / 8)
	// 突发码率比率: 1.35 (视情况取值, 大于等于1)
	durationPerPacket := float64(maxSize*1024*1024) / (float64(rate) * 1.35 / 8.0)
	wd, _ := os.Getwd()
	var index = 0
	wg := sync.WaitGroup{}
	sName := fmt.Sprintf("%x", md5.Sum([]byte(fileName)))
	var ret []string
	for st := 0.0; st < duration; st += durationPerPacket {
		wg.Add(1)
		index++
		filename := filepath.Join(wd, fmt.Sprint(sName, "_", index, ".mp4"))
		ret = append(ret, filename)
		go func(st float64, filename string) {
			defer wg.Done()
			// ffmpeg-ss 00:03:00 -i video.mp4 -t 60 -c copy -avoid_negative_ts 1 cut.mp4
			_ = exec.Command(
				"ffmpeg",
				"-ss", fmt.Sprint(st),
				"-i", fileName,
				"-t", fmt.Sprint(durationPerPacket),
				"-c", "copy",
				"-avoid_negative_ts", "1",
				filename,
			).Run()
		}(st, filename)
	}
	wg.Wait()
	return ret, nil
}
