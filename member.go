package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/line/line-bot-sdk-go/linebot"
)

var member map[string]int = map[string]int{} // 参加メンバー保存用 {player_name: game_num}

// スタンプからのメンバー登録
func registUser(bot *linebot.Client, event *linebot.Event) {
	user, _ := bot.GetProfile(event.Source.UserID).Do()

	if _, exist := member[user.DisplayName]; exist {
		newName := otherName(user.DisplayName)
		member[newName] = 0
		message := newName + "さんを追加しました"

		fmt.Println(message)

		if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message)).Do(); err != nil {
			log.Print(err)
		}
	} else {
		member[user.DisplayName] = 0
		message := user.DisplayName + "さんを追加しました"

		fmt.Println(message)

		if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message)).Do(); err != nil {
			log.Print(err)
		}
	}
}

// スタンプからのメンバー削除
func deleteUser(bot *linebot.Client, event *linebot.Event) {
	user, _ := bot.GetProfile(event.Source.UserID).Do()
	delete(member, user.DisplayName)
	message := user.DisplayName + "さんが退室しました"
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message)).Do(); err != nil {
		log.Print(err)
	}
}

// メッセージからのメンバー登録
func registGuest(bot *linebot.Client, name string, event *linebot.Event) {
	if _, exist := member[name]; exist {
		newName := otherName(name)
		member[newName] = 0
		message := newName + "さんを追加しました"

		fmt.Println(message)

		if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message)).Do(); err != nil {
			log.Print(err)
		}
	} else {
		member[name] = 0
		message := name + "さんを追加しました"

		fmt.Println(message)

		if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message)).Do(); err != nil {
			log.Print(err)
		}
	}
}

// メッセージからのメンバー削除
func deleteGuest(bot *linebot.Client, name string, event *linebot.Event) {
	delete(member, name)
	message := name + "さんが退室しました"
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message)).Do(); err != nil {
		log.Print(err)
	}
}

// 名前入力時、追加か削除かを選択する Template message を作成
func selectRegistOrDelete(bot *linebot.Client, message *linebot.TextMessage, event *linebot.Event) {
	registBtn := &linebot.PostbackTemplateAction{
		Label: "追加する",
		Data:  "action=regist&name=" + message.Text,
	}

	deleteBtn := &linebot.PostbackTemplateAction{
		Label: "退室させる",
		Data:  "action=delete&name=" + message.Text,
	}

	template := linebot.NewConfirmTemplate(message.Text+"さんを", registBtn, deleteBtn)
	confirmMessgage := linebot.NewTemplateMessage("this is confirm template", template)

	if _, err := bot.ReplyMessage(event.ReplyToken, confirmMessgage).Do(); err != nil {
		log.Print(err)
	}
}

// 同じ名前が入力された場合、改変した名前で登録する
func otherName(name string) string {
	flag := false
	i := 0
	newName := ""

	for {
		newName = name + strconv.Itoa(i)
		if _, exist := member[newName]; !exist {
			flag = true
		}

		if flag {
			break
		}

		i++
	}

	return newName
}
