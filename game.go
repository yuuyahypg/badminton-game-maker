package main

import (
	"fmt"
	"log"
	"math/rand"
	"sort"
	"strconv"

	"github.com/line/line-bot-sdk-go/linebot"
)

var gameNum int = 0         // コート数
var games []Game = []Game{} // マッチング済みの試合
var gameMessage string = "" // マッチング済みの試合確認用文字列

type Game struct {
	players [4]string
}

// コート数を選択させる Template message を作成
func startMakeGame(bot *linebot.Client, event *linebot.Event) {
	Btn3 := &linebot.PostbackTemplateAction{
		Label: "3",
		Data:  "action=start&num=3",
	}

	Btn4 := &linebot.PostbackTemplateAction{
		Label: "4",
		Data:  "action=start&num=4",
	}

	Btn5 := &linebot.PostbackTemplateAction{
		Label: "5",
		Data:  "action=start&num=5",
	}

	Btn6 := &linebot.PostbackTemplateAction{
		Label: "6",
		Data:  "action=start&num=6",
	}

	template := &linebot.ButtonsTemplate{
		Text:    "コート数を選択してください",
		Actions: []linebot.TemplateAction{Btn3, Btn4, Btn5, Btn6},
	}
	buttonMessgage := linebot.NewTemplateMessage("this is button template", template)

	if _, err := bot.ReplyMessage(event.ReplyToken, buttonMessgage).Do(); err != nil {
		log.Print(err)
	}
}

// 試合のマッチングを行う
// 各メンバーの試合数でソートして、少ない人から優先して試合を組む
func newGame(bot *linebot.Client, event *linebot.Event) {
	games = []Game{}
	if len(member) < (gameNum * 4) {
		gameNum = len(member) / 4
	}

	fmt.Println(gameNum)

	p := List{}
	for k, v := range member {
		e := Entry{k, v}
		p = append(p, e)
	}

	sort.Sort(p)

	players := p[:(gameNum * 4)]
	shuffle(players)

	for i := 0; i < gameNum; i++ {
		member[players[i*4].name] = member[players[i*4].name] + 1
		member[players[i*4+1].name] = member[players[i*4+1].name] + 1
		member[players[i*4+2].name] = member[players[i*4+2].name] + 1
		member[players[i*4+3].name] = member[players[i*4+3].name] + 1

		game := Game{
			players: [4]string{players[i*4].name, players[i*4+1].name, players[i*4+2].name, players[i*4+3].name},
		}
		games = append(games, game)
	}

	gameMessage = makeGameMessage()

	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(gameMessage)).Do(); err != nil {
		log.Print(err)
	}
}

// マッチングされた試合確認用文字列を作成
func makeGameMessage() string {
	message := ""
	for i, game := range games {
		message = message + "第" + strconv.Itoa(i+1) + "コート\n"
		message = message + game.players[0] + " & " + game.players[1] + "\n"
		message = message + game.players[2] + " & " + game.players[3] + "\n"

		if i < len(games)-1 {
			message = message + "\n\n"
		}
	}

	return message
}

func viewGames(bot *linebot.Client, event *linebot.Event) {
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(gameMessage)).Do(); err != nil {
		log.Print(err)
	}
}

// メンバーソート用
type Entry struct {
	name  string
	value int
}
type List []Entry

func (l List) Len() int {
	return len(l)
}

func (l List) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l List) Less(i, j int) bool {
	if l[i].value == l[j].value {
		return (l[i].name < l[j].name)
	} else {
		return (l[i].value < l[j].value)
	}
}

// メンバーシャッフル用
func shuffle(data []Entry) {
	n := len(data)
	for i := n - 1; i >= 0; i-- {
		j := rand.Intn(i + 1)
		data[i], data[j] = data[j], data[i]
	}
}
