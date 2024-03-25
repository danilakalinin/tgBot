package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/crocone/tg-bot"
	"github.com/joho/godotenv"
)

var bot *tgbotapi.BotAPI

type button struct {
	name string
	data string
}

func startMenu() tgbotapi.InlineKeyboardMarkup {
	states := []button{
		{
			name: "Забрать подарок 🎁",
			data: "get gift",
		},
		{
			name: "Начать психологический тест",
			data: "get test",
		},
	}

	buttons := make([][]tgbotapi.InlineKeyboardButton, len(states))
	for index, state := range states {
		buttons[index] = tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(state.name, state.data))
	}

	return tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(".env not loaded")
	}

	bot, err = tgbotapi.NewBotAPI(os.Getenv("TG_API_BOT_TOKEN"))
	if err != nil {
		log.Fatalf("Failed to initialize Telegram bot API: %v", err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Failed to start listening for updates %v", err)
	}

	for update := range updates {
		if update.CallbackQuery != nil {
			callbacks(update)
		} else if update.Message.IsCommand() {
			commands(update)
		} else {
			// simply message
		}
	}
}

func callbacks(update tgbotapi.Update) {
	data := update.CallbackQuery.Data
	chatId := update.CallbackQuery.Message.Chat.ID
	firstName := update.CallbackQuery.From.FirstName

	var text string
	switch data {

	case "get gift":
		text = fmt.Sprintf(`Здравствуйте, %v!

Вы попали в пространство семейных отношений, где найдете ответы на вопросы о том, как:

💬 вернуть близость в отношениях, которой так давно не хватало,
💬 легко проходить кризисы в ваших отношениях,
💬 вернуть мир, покой и взаимопонимание в семью,
💬 научиться поддерживать партнера,
💬 наладить отношения со своими детьми, 💬 вернуть доверительные отношения,
💬 выстраивать личные границы, понимая границы партнера

И наконец, как стать той счастливой и дружной семьей, о которой вы мечтаете.

%v, Я хочу и могу помочь Вам с этим 🙏

Я - Елена Плевако, семейный психолог. Помогаю женщинам 30+ построить счастливые отношения заново через раскрытие ценности и управление своим состоянием.
За 25 лет практики я помогла сохранить более 300 семей.

Для Вас я приготовила ПОДАРКИ:

🌺 Инструкция "5 простых шагов к обретению взаимопонимания с партнером за 7 дней, даже если без ссор не проходит и дня"
🌺 Моя сессия-разбор 🎁`, firstName, firstName)
	case "get test":
		msgConfig, inlineMarkup := testQuestions(chatId)
		msgConfig.ReplyMarkup = inlineMarkup
		sendMessage(msgConfig)

	default:
		text = "Неизвестная команда"
	}
	msg := tgbotapi.NewMessage(chatId, text)
	sendMessage(msg)
}

func commands(update tgbotapi.Update) {
	command := update.Message.Command()
	switch command {
	case "start":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Выберите действие")
		msg.ReplyMarkup = startMenu()
		msg.ParseMode = "Markdown"
		sendMessage(msg)
	}
}

func testQuestions(chatId int64) (tgbotapi.MessageConfig, tgbotapi.InlineKeyboardMarkup) {
	text := "Выберите один из вариантов ответа на вопрос:"

	// Создаем массив вопросов
	questions := []string{
		"Как часто вы чувствуете усталость?",
		// Добавьте здесь другие вопросы
	}

	// Создаем массив кнопок для каждого вопроса
	buttons := make([][]tgbotapi.InlineKeyboardButton, len(questions))
	for i, question := range questions {
		buttons[i] = make([]tgbotapi.InlineKeyboardButton, 3) // Предполагаем, что у нас три варианта ответа
		buttons[i][0] = tgbotapi.NewInlineKeyboardButtonData("Редко", "answer1_"+question)
		buttons[i][1] = tgbotapi.NewInlineKeyboardButtonData("Иногда", "answer2_"+question)
		buttons[i][2] = tgbotapi.NewInlineKeyboardButtonData("Часто", "answer3_"+question)
	}

	// Возвращаем текст и клавиатуру с кнопками
	return tgbotapi.NewMessage(chatId, text), tgbotapi.NewInlineKeyboardMarkup(buttons...)
}

func sendMessage(msg tgbotapi.Chattable) {
	if _, err := bot.Send(msg); err != nil {
		log.Panicf("Send message error: %v", err)
	}
}
