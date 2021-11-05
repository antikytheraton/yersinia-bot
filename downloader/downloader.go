package downloader

import (
	"bytes"
	"fmt"
	"net/url"
	"os/exec"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	command string = "youtube-dl -f bestvideo+bestaudio"
)

type YtDownloader interface {
	DownloadVideo(url string) error
}

type ytdl struct {
	command string
	timeout time.Duration
}

func (d *ytdl) execCommand(args ...string) (out string, err error) {
	var stderr, stdout bytes.Buffer

	cmd := exec.Command(d.command, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	time.AfterFunc(d.timeout, func() {
		if cmd.Process != nil && cmd.Process.Pid != 0 {
			out = out + fmt.Sprintf("\nexecute timeout: %v seconds.", d.timeout.Seconds())
			err = cmd.Process.Kill()
		}
	})
	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("%v\n%v", err, stderr.String())
	}
	out = stdout.String()
	return out, err
}

func (d *ytdl) DownloadVideo(rawURL string) error {
	videoURL, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	out, err := d.execCommand(videoURL.Host)
	log.Info(out)
	return nil
}

func NewYtDownloader() YtDownloader {
	return &ytdl{
		command: command,
		timeout: 60 * time.Minute,
	}
}
