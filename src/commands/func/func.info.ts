import type { AppFunc, BaseSession } from 'kbotify'
import { AppCommand, Card } from 'kbotify'

class InfoMsg extends AppCommand {
  code = 'info'
  trigger = 'info'
  help = '***info***'

  intro = 'info'
  func: AppFunc<BaseSession> = async (session) => {
    const kmd = `
        **欢迎来到星云工艺 NebulaeCraft**
        Wiki: [NebulaeWiki](https://wiki.knebulae.com/)
        BBS: [NebulaeBBS](https://bbs.knebulae.com/)
        ---
        机器人: [NebulaeRobot-Kaihaila](https://github.com/NebulaeCraft/NebulaeRobot-Kaihaila)
        欢迎Issue/PR
        `
    const cardObj = new Card({
      type: 'card',
      theme: 'secondary',
      size: 'lg',
      modules: [
        {
          type: 'section',
          text: {
            type: 'kmarkdown',
            content: kmd,
          },
        },
      ],
    })
    session.sendCard(cardObj)
  }
}

export const info = new InfoMsg()

