# gobut

## Configuration and running go-bot

1. `cp .env.sample .env`
1. `$EDITOR .env` and set your environment variables
1. `go run main.go`

### Environment variables

Variable | Description | Example
-------- | ----------- | -------
`DATABASE_URL` | Postgres database URL (optional) | postgres://user:passwd@host.com:5432/db
`BOTNAME` | The bot's IRC nickname | MyBot
`IRC_ADDRESS` | The address to the IRC server | irc.example.com:6667
`IRC_CHANNEL` | A comma-separated list of channels to join | #bots,#people
`POLL_USERNAMES` | Steam usernames to poll for game status | foxhop,japherwocky
`STEAMSTATUS_API` | URL to the Steam status API we use for following player activity | http://steamstatus-api.example:10000
`URINFO_API` | URL to the urinfo API we use for saving URLs | http://urinfo-api.example:10000/fetch
