import type { AppFunc, BaseSession } from 'kbotify'
import { AppCommand } from 'kbotify'
import config from '../config/config'
import { player } from './musicHandler/player'

class ChannelMenu extends AppCommand {
  code = 'channel'
  trigger = '频道'
  help = '点歌'

  intro = '点歌'
  func: AppFunc<BaseSession> = async (session) => {
    if (session.args.length === 0) {
      session.send(`
      请输入频道索引:
      1: 星云娘(默认)
      2: 直播区
      3: 游戏区
      4: 麻将区
      5: 聊天区
      `)
    }
    else {
      player.status.channel = config.channel[parseInt(session.args[0]) as keyof typeof config.channel]
      session.send(`已切换到频道${session.args[0]}`)
    }
  }
}

export const channel = new ChannelMenu()
