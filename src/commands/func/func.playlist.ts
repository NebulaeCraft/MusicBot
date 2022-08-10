import type { AppFunc, BaseSession } from 'kbotify'
import { AppCommand } from 'kbotify'
import { player } from './musicHandler/player'

class PlaylistMenu extends AppCommand {
  code = 'playlist'
  trigger = 'list'
  help = '点歌'

  intro = '点歌'
  func: AppFunc<BaseSession> = async (session) => {
    // session.send(session.args[0])
    player.status.playlist.forEach((item) => {
      session.send(`${item.name}--${item.ar}`)
    })
  }
}

export const playlist = new PlaylistMenu()

