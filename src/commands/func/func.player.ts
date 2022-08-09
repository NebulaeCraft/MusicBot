import type { AppFunc, BaseSession } from 'kbotify'
import { AppCommand, Card } from 'kbotify'
import { player } from './musicHandler/player'

class MusicMenu extends AppCommand {
  code = 'play'
  trigger = '点歌'
  help = '点歌'

  intro = '点歌'
  func: AppFunc<BaseSession> = async (session) => {
    // session.send(session.args[0])

    session.args.forEach(async (arg) => {
      let id: number
      if (arg.includes('music.163.com'))
        id = parseInt(arg.slice(arg.indexOf('?id=') + 4))

      else
        id = Number(arg)

      const testmusic = await player.fetch(id)
      console.log(JSON.stringify(testmusic))

      await player.push(testmusic).then(async () => {
        session.sendCard(new Card().addText(`${testmusic.name} 已加入播放列表`).toString())
        setTimeout(() => {
          if (!player.status.isPlaying)
            player.play(session)
        }, 5000)
      })
    })
  }
}

export const music = new MusicMenu()
