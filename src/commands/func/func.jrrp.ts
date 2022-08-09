import type { AppFunc, BaseSession } from 'kbotify'
import { AppCommand } from 'kbotify'

class JrrpMenu extends AppCommand {
  code = 'jrrp'
  trigger = '今日人品'
  help = '今日人品'

  intro = '今日人品'
  func: AppFunc<BaseSession> = async (session) => {
    session.send(() => (Math.random() * 100).toFixed(0).toString())
  }
}

export const jrrp = new JrrpMenu()

