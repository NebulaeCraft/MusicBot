import type { AppFunc, BaseSession } from 'kbotify'
import { AppCommand } from 'kbotify'
import axios from 'axios'
import { player } from './musicHandler/player'

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

            await player.push(testmusic).then(async () => {
              session.send(`${song.name} 已加入播放列表`)
              if (!player.status.isPlaying)
                player.play(session)
            })
          })
        })
    })
  }
}

export const list = new ListMenu()

