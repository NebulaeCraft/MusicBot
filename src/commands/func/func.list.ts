import type { AppFunc, BaseSession } from 'kbotify'
import { AppCommand, Card } from 'kbotify'
import axios from 'axios'
import { type Music, player } from './musicHandler/player'

class ListMenu extends AppCommand {
  code = 'play'
  trigger = '歌单'
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

      axios.get(`http://163.ishirai.cc:3000/playlist/detail?id=${id}`)
        .then(({ data }) => {
          const tracks = data.playlist.tracks

          tracks.forEach(async (song: any) => {
            const testmusic = await player.fetch(Number(song.id))
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
        })
    })
  }
}

export const list = new ListMenu()

