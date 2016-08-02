package main

import (
	"fmt"
	"log"

	trello "github.com/VojtechVitek/go-trello"
	"github.com/nlopes/slack"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("設定ファイル読み込みエラー: %s \n", err))
	}

	msg := fetchTrelloData()
	sendSlackMessage(msg)
}

func fetchTrelloData() string {
	trelloAppKey := viper.GetString("trello.appkey")
	trelloToken := viper.GetString("trello.token")
	trelloBoard := viper.GetString("trello.board")

	trello, err := trello.NewAuthClient(trelloAppKey, &trelloToken)
	if err != nil {
		log.Fatal(err)
	}

	board, err := trello.Board(trelloBoard)
	if err != nil {
		log.Fatal(err)
	}

	lists, err := board.Lists()
	if err != nil {
		log.Fatal(err)
	}

	msg := board.Name + "\n"

	for _, list := range lists {
		msg = msg + "*" + list.Name + "*" + "\n"

		cards, _ := list.Cards()
		for _, card := range cards {
			msg = msg + " - " + "_" + card.Name + "_" + "\n"
		}
	}

	return msg
}

func sendSlackMessage(msg string) {
	slackToken := viper.GetString("slack.token")
	slackChannel := viper.GetString("slack.channel")
	slackMessage := viper.GetString("slack.message")

	msg = slackMessage + "\n\n" + msg

	api := slack.New(slackToken)

	// 必ずボットをチャンネルに招待すること
	params := slack.PostMessageParameters{}
	params.AsUser = true
	params.LinkNames = 1

	api.PostMessage(slackChannel, msg, params)
}
