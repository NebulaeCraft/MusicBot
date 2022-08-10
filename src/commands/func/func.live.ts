import type { AppFunc, BaseSession } from 'kbotify'
import { AppCommand } from 'kbotify'

class AliveMenu extends AppCommand {
  code = 'alive'
  trigger = 'alive'
  help = 'alive'

  intro = 'alive'
  func: AppFunc<BaseSession> = async (session) => {
    const msg = await session.send(`开始保活 Start Time: ${new Date().toLocaleString()}`)
    function keepAlive() {
      setTimeout(() => {
        session.updateMessage(msg.msgSent!.msgId, `保活中 Update Time: ${new Date().toLocaleString()}`)
        console.log(`keepAlive: ${new Date().toLocaleString()}`)
        keepAlive()
      }, 15 * 60 * 1000)
    }
    keepAlive()
    console.log(`keepAlive: ${new Date().toLocaleString()}`)
  }
}

export const alive = new AliveMenu()

