# Trade Builder

Trade Builder Bot is a Discord bot designed for the ROBLOX game Bee Swarm Simulator. After trading was introduced to the game, players needed an easy and simple way to create and share their trades, leading to the development of Trade Builder Bot.

Trade Builder Bot is closely integrated with its [website](https://tradebuilder.app), and users can post trades to the website directly from Discord. The website is designed to be user-friendly and easy to navigate, with a simple and clean design.

## Demos

![Trade Builder Bot](https://cloud-punjwk6qj-hack-club-bot.vercel.app/0screenshot_2024-12-01_at_1.00.49_pm.png)

![Trade Builder Bot Video](https://cloud-g9js8bv0i-hack-club-bot.vercel.app/0screen_recording_2025-01-31_at_7.40.25___pm.mp4)

## Installation

Invite the bot here: https://trade.meta-bee.com/invite

Website: https://tradebuilder.app

You may view the website's Github repository here: https://github.com/alaninnovates/trade-builder-web

## Usage

### Trade commands

#### `/trade create`
Start by creating a trade with `/trade create`

#### `/trade lookingfor <sticker_name> <qty>`
#### `/trade offering <sticker_name> <qty>`
Add items to your trade with `/trade lookingfor` and `/trade offering`
- If the sticker is not in the trade, it will be added
- If the sticker is already in the trade, its quantity will be replaced (not added)

#### `/trade remove <type:lf/offering> <sticker_name>`
Remove items from your trade with `/trade remove`

#### `/trade view`
Displays your trade in an image that is ready to share.

At the bottom of the viewer, there are 3 buttons:

- Add looking for: Add a sticker to your looking for list
- Add offering: Add a sticker to your offering list
- Rerender: Rerender the trade viewer

#### `/trade info`

Returns info about the trade in textual format such as which stickers are in lf/offering. Essentially a copy pastable trade message from your trade

### Saving
#### `/trade save <save_name>`
Save your trade under a name
#### `/trade saves list`
List all your saved trades with their ids
#### `/trade saves load <id>`
Load a saved trade with a id, overwriting your current trade
#### `/trade saves delete <id>`
Delete a saved trade with a id

### Website commands

#### `/post <trade_id> [expire_time] [server_sync:true]`

Post a saved trade to the website
- trade_id: the id of your SAVED trade
- expire_time: how long before your trade post expires (in hh:mm:ss or mm/dd)
- server_sync: post to servers subscribed to server sync

### Market demand commands

#### `/top demand <duration> [category]`
#### `/top offer <duration> [category]`
View the top most demanded/offered stickers in the past day or week
- duration: day/week
- category: all/cubs/hives/bees/bears/mobs/critters/nectars/flowers/puffshrooms/leaves/tools/star signs/beesmas lights/field stamps

### Sync commands

Sync posts all new trades on the website to servers that are subscribed

**Premium server manager only:**

#### `/serversync setup <channel>`

Set up server sync in the channel - limit 1 channel per server

#### `/serversync remove`

Remove server sync for the entire server