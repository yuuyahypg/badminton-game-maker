# badminton-game-maker
## Requirement
* golang 1.8
* npm
* localtunnel

## Install
`npm install -g localtunnel`

## Getting Start
LINE Messaging APIのchannel secretとchannel access tokenを取得する
```
$ export BAD_SECRET=<channel secret>
$ export BAD_ACCESS_TOKEN=<channel access token>
```

公開するポートを登録  
`$ export BAD_PORT=<your_port>`

サーバーを起動する  
```
$ cd badminton-game-maker
$ go build
$ ./badminton-game-maker
```

localtunnelでポートを公開する  
`lt --port <your_port> --subdomain <your_webhook_url>`
