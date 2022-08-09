import { AppCommand, type AppFunc, type BaseSession, Card } from 'kbotify'

class StatusMenu extends AppCommand {
  code = 'status'
  trigger = '服务器状态'
  help = '服务器状态'

  intro = '服务器状态'
  func: AppFunc<BaseSession> = async (session) => {
    const card = new Card(`{
              "type": "card",
              "theme": "secondary",
              "size": "lg",
              "modules": [
                {
                  "type": "container",
                  "elements": [
                    {
                      "type": "image",
                      "src": "https://status.knebulae.com/?sid=1"
                    }
                  ]
                }
              ]
            }`)

    session.sendCard(card)
  }
}

export const status = new StatusMenu()

