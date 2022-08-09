import { bot } from './init/client'
import { jrrp } from './commands/func/func.jrrp'
import { play } from './commands/func/func.play'
import { music } from './commands/func/func.player'
import { list } from './commands/func/func.list'
import { skip } from './commands/func/func.skip'
import { channel } from './commands/func/func.channel'
import { lyric } from './commands/func/func.lyric'
import { volume } from './commands/func/func.volume'

bot.addCommands(jrrp)
bot.addCommands(play)
bot.addCommands(music)
bot.addAlias(music, 'd', 'dg')
bot.addCommands(channel)
bot.addAlias(channel, 'c', 'ch')
bot.addCommands(list)
bot.addCommands(skip)
bot.addCommands(lyric)
bot.addAlias(lyric, 'g', 'gc')
bot.addCommands(volume)
bot.addAlias(volume, 'v', 'y')
bot.connect()
