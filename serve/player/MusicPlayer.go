package player

import (
	"MusicBot/config"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gammazero/deque"
	"github.com/lonelyevil/kook"
)

type Music struct {
	ID       string
	Name     string
	Artists  []string
	Album    string
	File     string
	LastTime int
}

type ThreadManager struct {
	cmd    *exec.Cmd
	cancel func(*exec.Cmd)
	mu     sync.Mutex
	onExit func()
}

type Player struct {
	NowPlaying      *Music
	Queue           *deque.Deque[*Music]
	ControlChan     chan int
	Mutex           sync.Mutex
	Manager         *ThreadManager
	Ctx             *kook.KmarkdownMessageContext
	Channel         int
	Volume          int
	IsPaused        bool
	DefaultPlaylist *[]Music
}

var MusicPlayer *Player

func (music *Music) Copy() *Music {
	var artists []string
	for _, artist := range music.Artists {
		artists = append(artists, artist)
	}

	return &Music{
		ID:       music.ID,
		Name:     music.Name,
		Artists:  artists,
		Album:    music.Album,
		File:     music.File,
		LastTime: music.LastTime,
	}
}

func (tm *ThreadManager) StartThread(onExit func(), player *Player) error {
	logger := config.Logger

	tm.mu.Lock()
	defer tm.mu.Unlock()

	if tm.cmd != nil && tm.cancel != nil {
		tm.cancel(tm.cmd)
		tm.cmd.Wait()
	}

	cmd := exec.Command(config.Config.KOOKVoice,
		"-c",
		strconv.Itoa(player.Channel),
		"-i",
		player.NowPlaying.File,
		"-t",
		config.Config.BotToken,
		"-af",
		fmt.Sprintf("loudnorm=i=-16,volume=%ddB", player.Volume),
	)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	tm.cancel = func(cmd *exec.Cmd) {
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	}

	tm.cmd = cmd
	tm.onExit = onExit

	go func() {
		if err := cmd.Run(); err != nil {
			logger.Err(err).Msg("Error when running thread")
		}

		time.Sleep(5 * time.Second)

		if tm.onExit != nil {
			tm.onExit()
		}
	}()

	return nil
}

func (tm *ThreadManager) StopThread() {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if tm.cmd != nil && tm.cancel != nil {
		tm.cancel(tm.cmd)
		tm.cmd.Wait()
		tm.cmd = nil
		tm.cancel = nil
	}
}

func NewPlayer() *Player {
	player := new(Player)
	logger := config.Logger

	err := os.RemoveAll("./assets/music")
	if err != nil {
		logger.Err(err).Msg("Error when removing music folder")
	}
	err = os.Mkdir("./assets/music", 0755)
	if err != nil {
		logger.Err(err).Msg("Error when creating music folder")
	}

	player.Queue = deque.New[*Music]()
	player.ControlChan = make(chan int)
	player.Mutex = sync.Mutex{}
	player.Manager = &ThreadManager{}
	player.Channel = config.Config.VoiceChannel[0].ID
	player.Volume = config.Config.DefaultVolume
	player.IsPaused = true
	player.DefaultPlaylist = nil
	go player.Worker()

	return player
}

func (player *Player) AddMusic(music *Music) {
	player.AddMusicSilent(music)
	player.SendMsg(fmt.Sprintf("%s 已加入播放列表", music.Name))
}

func (player *Player) AddMusicSilent(music *Music) {
	player.Mutex.Lock()
	defer player.Mutex.Unlock()

	player.Queue.PushBack(music)
}

func (player *Player) RemoveMusic(musicID string) string {
	player.Mutex.Lock()
	defer player.Mutex.Unlock()

	for i := 0; i < player.Queue.Len(); i++ {
		music := player.Queue.At(i)
		if music.ID == musicID {
			name := music.Name
			player.Queue.Remove(i)
			return name
		}
	}

	return ""
}

func (player *Player) MoveUpMusic(musicID string) string {
	player.Mutex.Lock()
	defer player.Mutex.Unlock()

	for i := 0; i < player.Queue.Len(); i++ {
		music := player.Queue.At(i)
		if music.ID == musicID && i != 0 {
			player.Queue.Remove(i)
			player.Queue.Insert(i-1, music)
			return music.Name
		}
	}

	return ""
}

func (player *Player) MoveDownMusic(musicID string) string {
	player.Mutex.Lock()
	defer player.Mutex.Unlock()

	for i := 0; i < player.Queue.Len(); i++ {
		music := player.Queue.At(i)
		if music.ID == musicID && i != player.Queue.Len()-1 {
			player.Queue.Remove(i)
			player.Queue.Insert(i+1, music)
			return music.Name
		}
	}

	return ""
}

func (player *Player) MoveTopMusic(musicID string) string {
	player.Mutex.Lock()
	defer player.Mutex.Unlock()

	for i := 0; i < player.Queue.Len(); i++ {
		music := player.Queue.At(i)
		if music.ID == musicID {
			player.Queue.Remove(i)
			player.Queue.PushFront(music)
			return music.Name
		}
	}

	return ""
}

func (player *Player) SetCtx(ctx *kook.KmarkdownMessageContext) {
	player.Ctx = ctx
}

func (player *Player) SkipMusic() {
	player.Mutex.Lock()
	defer player.Mutex.Unlock()

	if player.NowPlaying != nil {
		player.SendMsg(fmt.Sprintf("即将跳过歌曲 %s", player.NowPlaying.Name))
		player.Manager.StopThread()
	} else {
		player.SendMsg("当前没有歌曲在播放")
	}
}

func (player *Player) SetVolume(volume int) {
	player.Mutex.Lock()
	defer player.Mutex.Unlock()

	player.Volume = volume

	if player.NowPlaying == nil {
		return
	}

	player.Queue.PushFront(player.NowPlaying.Copy())
	player.Manager.StopThread()
}

func (player *Player) SetChannel(channel int) {
	player.Mutex.Lock()
	defer player.Mutex.Unlock()

	player.Channel = channel

	if player.NowPlaying == nil {
		return
	}

	player.Queue.PushFront(player.NowPlaying.Copy())
	player.Manager.StopThread()
}

func (player *Player) Worker() {
	logger := config.Logger

	for {
		if !player.IsPaused {
			if player.NowPlaying == nil && player.Queue.Len() != 0 {
				player.Mutex.Lock()

				player.NowPlaying = player.Queue.PopFront()
				logger.Info().Msg("Now playing: " + player.NowPlaying.Name)
				logger.Info().Msg("Queue length: " + strconv.Itoa(player.Queue.Len()))

				player.Mutex.Unlock()

				logger.Info().Msg("Starting thread: " + player.NowPlaying.Name)
				player.Manager.StartThread(func() {
					player.NowPlaying = nil
					logger.Info().Msg("Set NowPlaying to nil")
				}, player)
				player.SendMusicCard(player.NowPlaying)
			} else {
				if player.DefaultPlaylist != nil && player.NowPlaying == nil && player.Queue.Len() == 0 {
					randomMusic := rand.Intn(len(*player.DefaultPlaylist))
					player.AddMusicSilent(&(*player.DefaultPlaylist)[randomMusic])
					logger.Info().Msg("Add music from default playlist")
				}
			}
		}

		time.Sleep(1 * time.Second) // prevent busy loop
	}
}

func (player *Player) Pause() {
	player.Mutex.Lock()
	defer player.Mutex.Unlock()

	player.IsPaused = true
	player.Manager.StopThread()
	player.SendMsg("已暂停")
}

func (player *Player) Resume() {
	player.Mutex.Lock()
	defer player.Mutex.Unlock()

	player.IsPaused = false
	player.SendMsg("已恢复")
}

func (player *Player) SendMsg(content string) {
	_, _ = player.Ctx.Session.MessageCreate(&kook.MessageCreate{
		MessageCreateBase: kook.MessageCreateBase{
			TargetID: player.Ctx.Common.TargetID,
			Content:  content,
		},
	})
}

func (player *Player) SendMusicCard(musicReq *Music) {
	logger := config.Logger
	cardMsg := kook.CardMessageCard{
		Theme: kook.CardThemeSuccess,
		Size:  kook.CardSizeLg,
	}
	section := kook.CardMessageSection{
		Mode: kook.CardMessageSectionModeRight,
		Text: kook.CardMessageElementKMarkdown{
			Content: "**歌曲：** " + musicReq.Name + "\n**歌手：** " + strings.Join(musicReq.Artists, ", ") + "\n**时长：** " + time.Duration(musicReq.LastTime*1000000).String(),
		},
	}
	cardMsg.AddModule(section.SetAccessory(&kook.CardMessageElementImage{
		Src:  musicReq.Album,
		Size: "lg",
	}))
	cardMsgCtx, err := cardMsg.MarshalJSON()
	if err != nil {
		logger.Error().Err(err).Msg("Marshal card message failed")
		return
	}
	cardMsgCtxStr := fmt.Sprintf("[%s]", cardMsgCtx)
	_, _ = player.Ctx.Session.MessageCreate(&kook.MessageCreate{
		MessageCreateBase: kook.MessageCreateBase{
			TargetID: player.Ctx.Common.TargetID,
			Content:  cardMsgCtxStr,
			Type:     kook.MessageTypeCard,
		},
	})
}

func (player *Player) SendMusicList() {
	player.Mutex.Lock()
	defer player.Mutex.Unlock()

	logger := config.Logger

	cardMsg := kook.CardMessageCard{
		Theme: kook.CardThemeSuccess,
		Size:  kook.CardSizeLg,
	}
	if player.Queue.Len() == 0 {
		section := kook.CardMessageSection{
			Mode: kook.CardMessageSectionModeRight,
			Text: kook.CardMessageElementKMarkdown{
				Content: "**歌曲列表为空**",
			},
		}
		cardMsg.AddModule(section.SetAccessory(&kook.CardMessageElementButton{
			Theme: kook.CardThemeInfo,
			Text:  "好了好了，我知道了",
			Click: string(kook.CardMessageElementButtonClickReturnVal),
			Value: "CONFIRM",
		}))
	} else {
		for i := 0; i < player.Queue.Len(); i++ {
			musicCtx := player.Queue.At(i)

			section := kook.CardMessageSection{
				Mode: kook.CardMessageSectionModeRight,
				Text: kook.CardMessageElementKMarkdown{
					Content: "**歌曲：** " + musicCtx.Name + "\n**歌手：** " + strings.Join(musicCtx.Artists, ", ") + "\n**时长：** " + time.Duration(musicCtx.LastTime*1000000).String(),
				},
			}

			cardMsg.AddModule(&section)

			upBtn := kook.CardMessageElementButton{
				Theme: kook.CardThemeDanger,
				Text:  "⬆️",
				Click: string(kook.CardMessageElementButtonClickReturnVal),
				Value: "UP" + musicCtx.ID,
			}

			downBtn := kook.CardMessageElementButton{
				Theme: kook.CardThemeDanger,
				Text:  "⬇️",
				Click: string(kook.CardMessageElementButtonClickReturnVal),
				Value: "DOWN" + musicCtx.ID,
			}

			topBtn := kook.CardMessageElementButton{
				Theme: kook.CardThemeDanger,
				Text:  "🔝",
				Click: string(kook.CardMessageElementButtonClickReturnVal),
				Value: "TOP" + musicCtx.ID,
			}

			delBtn := kook.CardMessageElementButton{
				Theme: kook.CardThemeDanger,
				Text:  "🗑️",
				Click: string(kook.CardMessageElementButtonClickReturnVal),
				Value: "DEL" + musicCtx.ID,
			}

			btnGroup := make(kook.CardMessageActionGroup, 0)
			btnGroup = append(btnGroup, upBtn, downBtn, topBtn, delBtn)

			cardMsg.AddModule(&btnGroup)
		}
	}
	cardMsgCtx, err := cardMsg.MarshalJSON()
	if err != nil {
		logger.Error().Err(err).Msg("Marshal card message failed")
		return
	}

	cardMsgCtxStr := fmt.Sprintf("[%s]", cardMsgCtx)
	_, _ = player.Ctx.Session.MessageCreate(&kook.MessageCreate{
		MessageCreateBase: kook.MessageCreateBase{
			TargetID: player.Ctx.Common.TargetID,
			Content:  cardMsgCtxStr,
			Type:     kook.MessageTypeCard,
		},
	})
}

func (player *Player) GetMusicByID(id string) *Music {
	player.Mutex.Lock()
	defer player.Mutex.Unlock()

	for i := 0; i < player.Queue.Len(); i++ {
		music := player.Queue.At(i)
		if music.ID == id {
			return music
		}
	}

	return nil
}
