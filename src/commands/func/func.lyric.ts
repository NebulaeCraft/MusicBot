import type { AppFunc, BaseSession } from 'kbotify'
import { AppCommand } from 'kbotify'
import { player } from './musicHandler/player'

class LyricMenu extends AppCommand {
  code = 'lyric'
  trigger = '歌词'
  help = '歌词'

  intro = '点歌'
  func: AppFunc<BaseSession> = async (session) => {
    if (session.args.length === 0) {
      if (player.status.ifLyric)
        session.send('歌词放送')
      else
        session.send('歌词关闭')
    }
    else {
      player.status.ifLyric = (session.args[0] === '1')
      if (player.status.ifLyric)
        session.send('歌词放送')
      else
        session.send('歌词关闭')
    }
  }
}

export const lyric = new LyricMenu()
