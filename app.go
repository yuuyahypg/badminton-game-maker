package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	bot, err := linebot.New(os.Getenv("ULIS_SECRET"), os.Getenv("ULIS_ACCESS_TOKEN"))
	if err != nil {
		panic(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}

		for _, event := range events {
			switch event.Type {
			case linebot.EventTypeMessage:
				switch message := event.Message.(type) {
				case *linebot.TextMessage: // 入力された文字列を名前としてメンバー登録する
					selectRegistOrDelete(bot, message, event)
				case *linebot.StickerMessage: // スタンプメッセージに対するアクション
					switch message.StickerID {
					case "4": // メンバー登録用(白人間　1行目:1列目のスタンプ)
						fmt.Println("メンバーの登録を開始します:" + event.Timestamp.String())
						registUser(bot, event)
					case "13": // メンバー削除用(白人間　1行目:2列目のスタンプ)
						deleteUser(bot, event)
					case "2": // コート数を指定してから試合マッチング(白人間　1行目:3列目のスタンプ)
						fmt.Println("コート数の登録を開始します:" + event.Timestamp.String())
						startMakeGame(bot, event)
					case "10": // 新しい試合のマッチング(白人間　1行目:4列目のスタンプ)
						fmt.Println("新しいゲームを作成します:" + event.Timestamp.String())
						newGame(bot, event)
					case "17": // 現在マッチングされている試合の表示(白人間　2行目:1列目のスタンプ)
						fmt.Println("現在のゲームを表示します:" + event.Timestamp.String())
						viewGames(bot, event)
					case "401": // 参加メンバーの確認(白人間　2行目:1列目のスタンプ)
						fmt.Println("現在のメンバーを表示します:" + event.Timestamp.String())
						reply := ""
						for player, _ := range member {
							reply = reply + player + "\n"
						}

						if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
							log.Print(err)
						}
					}
				}
			case linebot.EventTypePostback:
				checkPostbackData(bot, event)
			}
		}
	})

	if err := http.ListenAndServe(":3105", nil); err != nil {
		panic(err)
	}

	fmt.Println("server is running!")
}

// Template message からのリクエストを解析
func checkPostbackData(bot *linebot.Client, event *linebot.Event) {
	params := strings.Split(event.Postback.Data, "&")
	action := strings.Split(params[0], "=")

	switch action[1] {
	case "regist":
		name := strings.Split(params[1], "=")
		registGuest(bot, name[1], event)
	case "delete":
		name := strings.Split(params[1], "=")
		deleteGuest(bot, name[1], event)
	case "start":
		s := strings.Split(params[1], "=")
		num, _ := strconv.Atoi(s[1])
		gameNum = num
		newGame(bot, event)
	}
}
