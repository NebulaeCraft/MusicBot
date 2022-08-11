import type { AppFunc, BaseSession } from 'kbotify'
import { AppCommand, Card } from 'kbotify'
import { type Music, player } from './musicHandler/player'

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
      if (testmusic !== false) {
        await player.push(testmusic as Music).then(async () => {
          session.sendCard(new Card().addText(`${(testmusic as Music).name} 已加入播放列表`).toString())
          setTimeout(() => {
            if (!player.status.isPlaying)
              player.play(session)
          }, 5000)
        })
      }
      else {
        session.send('歌曲信息拉取失败')
      }
    })
  }
}

export const music = new MusicMenu()
