import child_process from 'child_process'
import fs from 'fs'
import axios from 'axios'
import request from 'request'
import { Lrc } from 'lrc-kit'
import type { BaseSession } from 'kbotify'
import { Card } from 'kbotify'
import path from '../../config/config'

// export namespace playerthis.Status {
//   export let isPlaying = false
//   export let playerPID: number | undefined = undefined
//   export let khlPID: number | undefined = undefined
//   export let playlist: Music[] = []
//   export let lastPlayTime: number | undefined = undefined
//   export const apiBase1 = 'http://localhost:3000'
//   export const apiBase2 = 'http://localhost:3000'
//   // export let apiBase: string = "http://163.ishirai.cc:3000";
//   export interface Music {
//     id: number
//     name: string
//     dt: number
//     cover: string
//     ar: string
//   }
// }

export interface Music {
  id: number
  name: string
  dt: number
  cover: string
  ar: string
}

class PlayerStatus {
  isPlaying = false
  playerPID: number | undefined
  khlPID: number | undefined
  playlist: Music[] = []
  lastPlayTime: number | undefined
  channel: number | string | undefined
  lastChannel: number | string | undefined
  ifLyric = false
  volume = -25
  apiBase1 = 'http://163.ishirai.cc:3000'
  apiBase2 = 'http://163.ishirai.cc:3000'
  // apiBase1 = 'http://localhost:3000'
  // apiBase2 = 'http://localhost:3000'
  // export let apiBase: string = "http://163.ishirai.cc:3000";
}

class Player {
  status: PlayerStatus

  constructor() {
    this.status = new PlayerStatus()
    this.status.channel = 5892238130033188
    this.status.playlist = []
  }

  async downloadFile(uri: string, filename: string) {
    if (uri === undefined)
      throw new Error('uri undefined')
    return new Promise((resolve) => {
      const stream = fs.createWriteStream(filename)
      console.log(`start downloading ${filename}`)
      request(uri).pipe(stream).on('close', () => {
        console.log(`${filename} downloaded`)
        resolve(true)
      })
    })
  }

  async stop() {
    console.log('stopfunc')

    return new Promise((resolve) => {
      if (this.status.isPlaying) {
        if (this.status.playerPID !== undefined) {
          // child_process.exec(`taskkill /pid ${this.status.playerPID} -t -f`) // TODO
          child_process.exec(`kill ${this.status.khlPID}`) // TODO
          console.log(`killed ${this.status.playerPID}`)
        }
        else {
          console.log('pid undefined')
        }
      }
      resolve(true)
    })
  }

  async push(music: Music) {
    console.log(`push ${music.name}`)

    let url: string, filename: string, isSuccess: boolean
    await (axios.get(`${this.status.apiBase1}/song/url?id=${music.id}`, {
      headers: {
        Cookie: path.cookie,
      },
    })
      .then(({ data }) => {
        url = data.data[0].url
        filename = `./res/${music.id}.mp3`
        isSuccess = true
      }))
      .then(async () => {
        await this.downloadFile(url, filename)
          .catch((e: Error) => {
            console.log(`${music.name} download failed`)
            console.log(e)
            isSuccess = false
            return false
          })
          .then(() => {
            if (isSuccess)
              this.status.playlist.push(music)
          }).then(() => {
            return true
          })
      }).then(() => {
        return true
      }).then(() => {
        return isSuccess
      })
  }

  async play(session: BaseSession) {
    const music = this.status.playlist.shift()!

    console.log(`playing ${music.name}`)

    const lrcObj = await this.fetchlrc(music.id)
    this.sendInfoCard(music, session)
    if (this.status.lastPlayTime === undefined || (new Date().getTime() - this.status.lastPlayTime) > 3 * 60 * 1000 || String(this.status.channel) !== String(this.status.lastChannel)) {
      if (this.status.khlPID !== undefined) {
        // child_process.exec(`taskkill /pid ${this.status.khlPID} -t -f`) // TODO
        child_process.exec(`kill ${this.status.khlPID}`) // TODO
        console.log(`killed khl-voice ${this.status.khlPID}`)
      }
      this.status.lastChannel = this.status.channel
      const khlProcess = child_process.exec(`${path.khlvoice} -i zmq:tcp://localhost:5559 -t 1/MTEwNDc=/o6EyPxPpLSN4JXlvIuovpA== -c ${this.status.channel}`)// 星云娘
      this.status.khlPID = khlProcess.pid
      if (this.status.ifLyric)
        setTimeout(() => { this.sendLyric(lrcObj, session) }, 3000)
    }
    else {
      if (this.status.ifLyric)
        setTimeout(() => { this.sendLyric(lrcObj, session) }, 0)
    }
    const playerProcess = child_process.exec(`ffmpeg -re -i ./res/${music.id}.mp3 -af volume=${this.status.volume}dB -ab 120k -acodec libopus -f mpegts zmq:tcp://127.0.0.1:5559`)
    this.status.playerPID = playerProcess.pid
    this.status.isPlaying = true
    this.status.lastPlayTime = new Date().getTime() + music.dt + 2000
    // return true;
    return new Promise((resolve) => {
      setTimeout(() => {
        if (this.status.playlist.length > 0) {
          this.play(session)
        }
        else {
          session.send('播放列表已完成')
          this.status.isPlaying = false
        }
        fs.rmSync(`./res/${music.id}.mp3`)
        resolve(music.name)
      }, music.dt + 2000)
    })
  }

  async fetch(musicID: number) {
    let music: Music
    return axios.get(`${this.status.apiBase2}/song/detail?ids=${musicID}`, {
      headers: {
        Cookie: path.cookie,
      },
    })
      .then(({ data }) => {
        const arArray: string[] = []
        data.songs[0].ar.forEach((ar: any) => {
          arArray.push(ar.name)
        })
        music = { id: musicID, name: data.songs[0].name, dt: data.songs[0].dt, cover: `${data.songs[0].al.picUrl}?param=130y130`, ar: arArray!.toString() }
        return music
      }).catch((e: Error) => {
        console.log(e)
        return false
      })
  }

  async fetchlrc(musicID: number) {
    return axios.get(`${this.status.apiBase2}/lyric?id=${musicID}`, {
      headers: {
        Cookie: path.cookie,
      },
    })
      .then(({ data }) => {
        if (data?.romalrc?.lyric) {
          const lrcobj = { lyric: data.lrc.lyric, romalrc: data?.romalrc?.lyric }
          return lrcobj
        }
        else {
          const lrcobj = { lyric: data.lrc.lyric }
          return lrcobj
        }
      })
  }

  async sendLyric(lrcObj: any, session: BaseSession) {
    if (lrcObj?.romalrc) {
      console.log('send lyric1')
      const lrcLyric = Lrc.parse(lrcObj.lyric)
      const msgLyric = await session.sendCard(new Card().addText('歌词').toString())

      const lrcRomal = Lrc.parse(lrcObj.romalrc)
      const msgRomal = await session.sendCard(new Card().addText('罗马音').toString())

      lrcLyric.lyrics.forEach(async (lyric) => {
        setTimeout(() => {
          try {
            session.updateMessage(msgLyric.msgSent!.msgId, new Card().addText(lyric.content).toString())
          }
          catch (error) {
            console.log(error)
          }
        }, Math.ceil(lyric.timestamp * 1000))
      })
      lrcRomal.lyrics.forEach(async (lyric) => {
        setTimeout(() => {
          try {
            session.updateMessage(msgRomal.msgSent!.msgId, new Card().addText(lyric.content).toString())
          }
          catch (error) {
            console.log(error)
          }
        }, Math.ceil(lyric.timestamp * 1000))
      })
    }
    else if (lrcObj?.lyric) {
      console.log('send lyric2')
      const lrcLyric = Lrc.parse(lrcObj.lyric)
      const msgLyric = await session.sendCard(new Card().addText('歌词').toString())

      lrcLyric.lyrics.forEach(async (lyric) => {
        setTimeout(() => {
          try {
            session.updateMessage(msgLyric.msgSent!.msgId, new Card().addText(lyric.content).toString())
          }
          catch (error) {
            console.log(error)
          }
        }, Math.ceil(lyric.timestamp * 1000))
      })
    }
  }

  sendInfoCard(music: Music, session: BaseSession) {
    const now = new Date()
    const start = new Date(now.getFullYear(), now.getMonth(), now.getDate())
    const end = new Date(start.getTime() + music.dt)
    const dnow = new Date(now.getTime() + music.dt)

    console.log(`splay:${music.name}`)
    console.log(`cov:${music.cover}`)

    session.sendCard(new Card({
      type: 'card',
      theme: 'success',
      size: 'lg',
      modules: [
        {
          type: 'section',
          text: {
            type: 'kmarkdown',
            content: `**歌曲：** ${music.name}\n**时长：**${`${end.getMinutes()}:${end.getSeconds()}`}\n**歌手：** ${music.ar}\n**开始播放：**  ${now.toString()} `,
          },
          mode: 'right',
          accessory: {
            type: 'image',
            src: music.cover,
            size: 'lg',
          },
        },
        {
          type: 'countdown',
          mode: 'second',
          startTime: `${now.getTime()}`,
          endTime: `${dnow.getTime()}`,
        },
      ],
    }))
  }
}

export const player = new Player()
