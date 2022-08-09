import type { AppFunc, BaseSession } from 'kbotify'
import { AppCommand } from 'kbotify'
import { player } from './musicHandler/player'

class VolumeMenu extends AppCommand {
  code = 'volume'
  trigger = '音量'
  help = '歌词'

  intro = '点歌'
  func: AppFunc<BaseSession> = async (session) => {
    if (session.args.length === 0) {
      session.send(`当前音量: ${player.status.volume}dB, 默认:-25dB`)
    }
    else {
      player.status.volume = parseInt(session.args[0])
      session.send(`已设置音量为${player.status.volume}dB`)
    }
  }
}

export const volume = new VolumeMenu()
