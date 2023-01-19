# golang_tg_saver2saved

A simple Telegram bot to send photos and videos from a directory to all users. 

Initially created mostly to automate converting videos so that they show up nicely in the messengers, but hey! why not send them automatically as well.

I tried to send things to Saved messages in Telegram, but the API is quite complicated and the Bot implementation turned out to be so much easier. Currently I even use direct https API calls (wanted to see how they work in Go)

# Install/run

- `git clone https://github.com/eugenpt/golang_tg_saver2saved.git`
- `cd golang_tg_saver2saved`
- `go build`
- Get a bot from [@BotFather](https://t.me/BotFather) 
	- save token into `token.txt`
	- Write `/start` to the bot) so that it knows to whom to send stuff
- Save path to the desired dir to watch into `dir.txt`
- Run `golang_tg_saver2saved`


requires ffmpeg to convert videos
