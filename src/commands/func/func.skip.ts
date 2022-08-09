import type { AppFunc, BaseSession } from 'kbotify'
import { AppCommand } from 'kbotify'
import { player } from './musicHandler/player'

class SkipMenu extends AppCommand {
  code = 'skip'
  trigger = '切歌'
  help = '点歌'

  intro = '点歌'
  func: AppFunc<BaseSession> = async (session) => {
    if (player.status.playlist.length > 0)
      await player.stop().then(() => { player.play(session) })

    else
      session.send('播放列表为空')
  }
}

export const skip = new SkipMenu()
