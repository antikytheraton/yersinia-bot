package bot

import (
	"flag"
	"io/ioutil"
	"time"

	"github.com/NicoNex/echotron/v3"
	"github.com/antikytheraton/yersinia-bot/downloader"
	"github.com/peterbourgon/ff/v3"
	log "github.com/sirupsen/logrus"
)

var (
	fs    = flag.NewFlagSet("telegram-bot", flag.ExitOnError)
	token = fs.String("token", "", "telegram api token")
)

var dsp echotron.Dispatcher

type Bot interface {
	Run(downloader downloader.YtDownloader, args []string) error
}

type bot struct {
	chatID     int64
	api        echotron.API
	downloader downloader.YtDownloader
}

var _ echotron.Bot = (*bot)(nil)

func (b *bot) selfDestruct(timech <-chan time.Time) {
	<-timech
	b.api.SendMessage("goodbye", b.chatID, nil)
	dsp.DelSession(b.chatID)
}

func (b *bot) Update(update *echotron.Update) {
	var err error
	defer func() {
		if err != nil {
			log.Error(err)
		}
	}()

	switch update.Message.Text {
	case "/start":
		log.Info("/start")
		_, err = b.api.SendMessage("Hello world", b.chatID, nil)
		return

	case "/video":
		log.Info("/video")
		var cnt []byte
		cnt, err = ioutil.ReadFile("testdata/john_china.mp4")
		file := echotron.NewInputFileBytes("john_china.mp4", cnt)
		_, err = b.api.SendVideo(file, b.chatID, nil)
		return

	default:
		log.Info("/dice")
		_, err = b.api.SendDice(b.chatID, echotron.Die, nil)
		return
	}
}

func (b *bot) printStart() {
	info, err := b.api.GetMe()
	if err != nil {
		log.Error(err)
	}
	botName := info.Result.Username
	log.Infof("starting %s...", botName)
}

func new(token string, downloader downloader.YtDownloader) echotron.NewBotFn {
	return func(chatID int64) echotron.Bot {
		b := &bot{
			chatID:     chatID,
			api:        echotron.NewAPI(token),
			downloader: downloader,
		}
		go b.selfDestruct(time.After(5 * time.Minute))

		b.printStart()
		return b
	}
}

func Run(downloader downloader.YtDownloader, args []string) error {
	err := ff.Parse(fs, args, ff.WithEnvVarPrefix("TELEGRAM"))
	if err != nil {
		return err
	}
	dsp := echotron.NewDispatcher(*token, new(*token, downloader))
	return dsp.Poll()
}
