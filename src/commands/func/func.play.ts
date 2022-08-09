import type { AppFunc, BaseSession } from 'kbotify'
import { AppCommand } from 'kbotify'
import { player } from './musicHandler/player'

class PlayMenu extends AppCommand {
  code = 'play'
  trigger = '播放'
  help = '点歌'

  intro = '点歌'
  func: AppFunc<BaseSession> = async (session) => {
    player.play(session)
  }
}

export const play = new PlayMenu()
