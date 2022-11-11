package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jackc/pgx/v5"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"log"
	"strconv"
	"time"
)

func telega_chat_buttons() {

}

func telega_ban_check(Telegaid int64) bool {

	//DATABASE_URL := "postgres://gpsdata:working123@194.87.234.209/gpsdata"
	DATABASE_URL := "postgres://gpsdata:working123@194.87.234.209/registers"
	conn, err := pgx.Connect(context.Background(), DATABASE_URL)
	if err != nil {
		log.Printf("Unable to connect to database: %v\n", err)
		//os.Exit(1)
	}
	defer conn.Close(context.Background())

	//fmt.Println("telagaiid=", Telegaid)

	var id_telegram int

	rrrr := "select COUNT(*) from tgid_ban_id where id_telegram=$1"
	//SELECT COUNT(*) FROM tgid_ban_id where id_telegram = " + strconv.FormatInt(Telegramid, 10)
	err = conn.QueryRow(context.Background(), rrrr, Telegaid).Scan(&id_telegram)
	if err != nil {
		log.Printf("QueryRow failed: %v\n", err)
		//os.Exit(1)
	}

	if id_telegram > 0 {
		return true
		fmt.Printf("true=", id_telegram)

	}
	fmt.Printf("false", id_telegram)
	return false

}

func initiLogger() {
	path := "telegram-log"
	var writer, err = rotatelogs.New(
		fmt.Sprintf("%s-%s", path, "%Y-%m-%d"),
		rotatelogs.WithMaxAge(time.Hour*24),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		log.Fatalf("Failed to Initialize Log File %s", err)
	}

	log.SetFlags(5)
	log.SetOutput(writer)

	return
}

func main() {
	initiLogger()
	log.Println("Starting service !!!")
	bot, err := tgbotapi.NewBotAPI("1369821650:AAF4OI72Ncb-Q0rCO9YgtHv-vQjpLS02E3o")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			//ОСНОВНОЙ ЦИКЛ begin
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			//msg.ReplyToMessageID = update.Message.MessageID

			Telegaid := update.Message.Chat.ID
			//TeleText := update.Message.Text

			bot.Send(tgbotapi.NewMessage(Telegaid, "Ваш ID: "+strconv.FormatInt(update.Message.Chat.ID, 10)))
			log.Printf("ID: " + strconv.FormatInt(update.Message.Chat.ID, 10))

			//bot.Send(tgbotapi.NewMessage(Telegaid, "Ваш ТЕКСТ: "+update.Message.Text))
			//log.Printf("Ваш ТЕКСТ: " + update.Message.Text)

			if telega_ban_check(Telegaid) {
				var NumericKeyboard = tgbotapi.NewReplyKeyboard(
					tgbotapi.NewKeyboardButtonRow(
						tgbotapi.NewKeyboardButton("ТЕХПОДДЕРЖКА"),
					),
				)
				msg := tgbotapi.NewMessage(Telegaid, "Обнаружен в блок списке, обратитесь в техподдержку.")
				msg.ReplyMarkup = NumericKeyboard
				bot.Send(msg)

				//ОСНОВНОЙ ЦИКЛ  end
			}
		}
	}
}
