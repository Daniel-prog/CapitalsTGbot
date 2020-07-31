/* Этот простой Телеграм-бот выводит столицу любой страны мира при отправке ему названия этой страны.
В этом проекте не использовались готовые библиотеки на Go для Telegram.
*/

package main

import ( 
	"fmt"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)
/* Задаём структуры для json-ответа*/
type UpdateT struct {
	Ok bool `json:"ok"`
	Result []UpdateResultT `json:"result"`
}

type UpdateResultT struct {
	UpdateId int `json:"update_id"`
	Message UpdateResultMessageT `json:"message"`
}

type UpdateResultMessageT struct {
	MessageId int `json:"message_id"`
	From UpdateResultFromT `json:"from"`
	Chat UpdateResultChatT `json:"chat"`
	Date int `json:"date"`
	Text string `json:"text"`
	// Entities []UpdateResultEntitiesT `json:"entities, omitempty"`
}

type UpdateResultFromT struct {
	Id int `json:"int"`
	IsBot bool `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Username string `json:"username"`
	Language string `json:"language_code"`
}

type UpdateResultChatT struct {
	Id int `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Username string `json:"username"`
	Type string `json:"private"`
}

// type UpdateResultEntitiesT struct {
// 	Offset int `json:"offset"`
// 	Length int `json:"length"`
// 	Type string `json:"type"`

// }

type SendMessageResponseT struct {
	Ok bool `json:"ok"`
	Result UpdateResultMessageT `json:"result"`
}

const BaseTelegramUrl = "https://api.telegram.org"
const telegramToken = "1367052814:AAFtpyce_JBbFc1a58Rem8e8BsTJwbglDTY"
const getUpdatesUri = "getUpdates"
const sendMessageUrl = "sendMessage"

const keywordStart = "/start"

func main() {
	
	Capitals := map[string]string {
		"Австралия"	:	"Канберра",
		"Австрия":		"Вена",
		"Азербайджан":	"Баку",
		"Албания":		"Тирана",
		"Алжир"	:		"Алжир",
		"Ангола":		"Луанда",
		"Андорра":		"Андорра-ла-Велья",
		"Антигуа и Барбуда":	"Сент-Джонс",
		"Аргентина"	:	"Буэнос-Айрес",
		"Армения"	:	"Ереван",
		"Афганистан":	"Кабул",
		"Багамы"	:	"Нассау",
		"Бангладеш"	:	"Дакка",
		"Барбадос"	:	"Бриджтаун",
		"Бахрейн"	:	"Манама",
		"Беларусь"	:	"Минск",
		"Белиз"		:	"Бельмопан",
		"Бельгия"	:	"Брюссель",
		"Бенин"		:	"Порто-Ново",
		"Болгария"	:	"София",
		"Боливия"	:	"Сукре",
		"Босния и Герцеговина":	"Сараево",
		"Ботсвана"	:	"Габороне",
		"Бразилия"	:	"Бразилиа",
		"Бруней"	:	"Бандар-Сери-Багаван",
		"Буркина Фасо":	"Уагадугу",
		"Бурунди"	:	"Бужумбура",
		"Бутан"		:	"Тхимпху",
		"Вануату"	:	"Порт-Вила",
		"Ватикан"	:	"Ватикан",
		"Великобритания":	"Лондон",
		"Венгрия"	:	"Будапешт",
		"Венесуэла"	:	"Каракас",
		"Восточный" :	"Тимор	Дили",
		"Вьетнам"	:	"Ханой",
		"Габон"		:	"Либревиль",
		"Гаити"		:	"Порт-о-Пренс",
		"Гайана"	:	"Джорджтаун",
		"Гамбия"	:	"Банжул",
		"Гана"		:	"Аккра",
		"Гватемала"	:	"Гватемала",
		"Гвинея"	:	"Конакри",
		"Гвинея-Бисау":	"Бисау",
		"Германия"	:	"Берлин",
		"Гондурас"	:	"Тегусигальпа",
		"Гренада"	:	"Сент-Джорджес",
		"Греция"	:	"Афины",
		"Грузия"	:	"Тбилиси",
		"Дания"		:	"Копенгаген",
		"Джибути"	:	"Джибути",
		"Доминика"	:	"Розо",
		"Доминиканская Республика":	"Санто-Доминго",
		"Египет"	:	"Каир",
		"Замбия"	:	"Лусака",
		"Зимбабве"	:	"Хараре",
		"Израиль"	:	"Иерусалим",
		"Индия"		:	"Нью-Дели",
		"Индонезия"	:	"Джакарта",
		"Иордания"	:	"Амман",
		"Ирак"		:	"Багдад",
		"Иран"		:	"Тегеран",
		"Ирландия"	:	"Дублин",
		"Исландия"	:	"Рейкьявик",
		"Испания"	:	"Мадрид",
		"Италия"	:	"Рим",
		"Йемен"		:	"Сана",
		"Кабо-Верде":	"Прая",
		"Казахстан"	:	"Астана",
		"Камбоджа"	:	"Пномпень",
		"Камерун"	:	"Яунде",
		"Канада"	:	"Оттава",
		"Катар"		:	"Доха",
		"Кения"		:	"Найроби",
		"Кипр"		:	"Никосия",
		"Киргизия"	:	"Бишкек",
		"Кирибати"	:	"Южная Тарава",
		"Китай"		:	"Пекин",
		"Колумбия"	:	"Санта-Фе-де-Богота",
		"Коморы"	:	"Морони",
		"Демократическая Республика Конго":	"Киншаса",
		"Конго, демократическая республика":	"Киншаса",
		"Конго, республика":	"Браззавиль",
		"Конго": 		"Браззавиль",
		"Коста-Рика":	"Сан-Хосе",
		"Кот-д’Ивуар":	"Ямусукро",
		"Куба"		:	"Гавана",
		"Кувейт"	:	"Эль-Кувейт",
		"Лаос"		:	"Вьентьян",
		"Латвия"	:	"Рига",
		"Лесото"	:	"Масеру",
		"Либерия"	:	"Монровия",
		"Ливан"		:	"Бейрут",
		"Ливия"		:	"Триполи",
		"Литва"		:	"Вильнюс",
		"Лихтенштейн":	"Вадуц",
		"Люксембург":	"Люксембург",
		"Маврикий"	:	"Порт-Луи",
		"Мавритания":	"Нуакшот",
		"Мадагаскар":	"Антананариву",
		"Македония"	:	"Скопье",
		"Малави"	:	"Лилонгве",
		"Малайзия"	:	"Куала-Лумпур",
		"Мали"		:	"Бамако",
		"Мальдивы"	:	"Мале",
		"Мальта"	:	"Валлетта",
		"Марокко"	:	"Рабат",
		"Маршалловы Острова":	"Маджуро",
		"Мексика"	:	"Мехико",
		"Мозамбик"	:	"Мапуту",
		"Молдавия"	:	"Кишинев",
		"Монако"	:	"Монако",
		"Монголия"	:	"Улан-Батор",
		"Мьянма"	:	"Найпьидо",
		"Намибия"	:	"Виндхук",
		"Науру"		:	"официальной столицы не имеет",
		"Непал"		:	"Катманду",
		"Нигер"		:	"Ниамей",
		"Нигерия"	:	"Абуджа",
		"Нидерланды":	"Амстердам",
		"Никарагуа"	:	"Манагуа",
		"Новая Зеландия":	"Веллингтон",
		"Норвегия"	:	"Осло",
		"Арабские Эмираты":	"Абу-Даби",
		"Объединенные Арабские Эмираты":	"Абу-Даби",
		"Оман"		:	"Маскат",
		"Пакистан"	:	"Исламабад",
		"Палау"		:	"Мелекеок",
		"Панама"	:	"Панама",
		"Папуа"		:	"Порт-Морсби",
		"Папуа - Новая Гвинея":	"Порт-Морсби",
		"Парагвай"	:	"Асунсьон",
		"Перу"		:	"Лима",
		"Польша"	:	"Варшава",
		"Португалия":	"Лиссабон",
		"Россия"	:	"Москва",
		"Руанда"	:	"Кигали",
		"Румыния"	:	"Бухарест",
		"Сальвадор"	:	"Сан-Сальвадор",
		"Самоа"		:	"Апиа",
		"Сан-Марино":	"Сан-Марино",
		"Сан-Томе и Принсипи":	"Сан-Томе",
		"Саудовская Аравия":	"Эр-Рияд",
		"Свазиленд"	:	"Мбабане",
		"Северная Корея":	"Пхеньян",
		"Сейшелы"	:	"Виктория",
		"Сенегал"	:	"Дакар",
		"Сент-Винсент и Гренадины":	"Кингстаун",
		"Сент-Китс и Невис":	"Бастер",
		"Сент-Люсия":	"Кастри",
		"Сербия"	:	"Белград",
		"Сингапур"	:	"Сингапур",
		"Сирия"		:	"Дамаск",
		"Словакия"	:	"Братислава",
		"Словения"	:	"Любляна",
		"США"		:	"Вашингтон",
		"Соединенные Штаты Америки"	:"Вашингтон",
		"Соломоновы Острова":	"Хониара",
		"Сомали"	:	"Могадишо",
		"Судан"		:	"Хартум",
		"Суринам"	:	"Парамарибо",
		"Сьерра-Леоне":	"Фритаун",
		"Таджикистан":	"Душанбе",
		"Таиланд"	:	"Бангкок",
		"Танзания"	:	"Додома",
		"Того"		:	"Ломе",
		"Тонга"		:	"Нукуалофа",
		"Тринидад и Тобаго":	"Порт-оф-Спейн",
		"Тувалу"	:	"Фунафути",
		"Тунис"		:	"Тунис",
		"Туркмения"	:	"Ашхабад",
		"Турция"	:	"Анкара",
		"Уганда"	:	"Кампала",
		"Узбекистан":	"Ташкент",
		"Украина"	:	"Киев",
		"Уругвай"	:	"Монтевидео",
		"Федеративные штаты Микронезии"	:"Паликир",
		"Микронезия"	:"Паликир",
		"Фиджи"		:	"Сува",
		"Филиппины"	:	"Манила",
		"Финляндия"	:	"Хельсинки",
		"Франция"	:	"Париж",
		"Хорватия"	:	"Загреб",
		"Центрально-Африканская Республика":	"Банги",
		"ЦАР":	"Банги",
		"Чад"		:	"Нджамена",
		"Черногория":	"Подгорица",
		"Чехия"		:	"Прага",
		"Чили"		:	"Сантьяго",
		"Швейцария"	:	"Берн",
		"Швеция"	:	"Стокгольм",
		"Шри-Ланка"	:	"Коломбо",
		"Эквадор"	:	"Кито",
		"Экваториальная Гвинея":	"Малабо",
		"Эритрея"	:	"Асмэра",
		"Эстония"	:	"Таллин",
		"Эфиопия"	:	"Аддис-Абеба",
		"Южная Корея":	"Сеул",
		"Южно-Африканская Республика":	"Претория",
		"ЮАР":	"Претория",
		"Ямайка"	:	"Кингстон",
		"Япония"	:	"Токио",
	}

	updateNum := "1"

	for {

		update, err := getUpdates(updateNum)

		if err != nil {
			fmt.Println(err.Error())

			continue
		}

		for _, item := range(update.Result) {
			msgHello := "Привет, " + item.Message.From.FirstName + " " + item.Message.From.LastName + "!"
			msgHow := "Дела у меня отлично!"
			msgWhat := "Да так, столицы стран повторяю..."
			text := item.Message.Text

			if text == keywordStart {
				sendMessage(item.Message.Chat.Id, "Введите название любой страны мира (с большой буквы), а я назову её столицу \xF0\x9F\x8C\x8E \xF0\x9F\x8C\x8D \xF0\x9F\x8C\x8F")
				updateNum = strconv.Itoa(item.UpdateId + 1)
				continue
			}

			switch (true) {
				case strings.Contains(text, "Привет"), strings.Contains(text, "привет"):
					sendMessage(item.Message.Chat.Id, msgHello)

				case strings.Contains(text, "Как дела"), strings.Contains(text, "как дела"):
					sendMessage(item.Message.Chat.Id, msgHow)

				case strings.Contains(text, "Что делаеш"), strings.Contains(text, "что делаеш"):
					sendMessage(item.Message.Chat.Id, msgWhat)

				default:
					if Capitals[text] != "" {
						msg := "Столица страны " + text + " - " + Capitals[text] + "!"
						sendMessage(item.Message.Chat.Id, msg)
					} else {
						sendMessage(item.Message.Chat.Id, "Не могу разобрать, повторите запрос")
					}
			}

			updateNum = strconv.Itoa(item.UpdateId + 1)
		}
	}
}

func getUpdates(num string) (UpdateT, error){
	url := BaseTelegramUrl + "/bot" + telegramToken + "/" + getUpdatesUri + "?offset=" + num
	response := getResponse(url)

	update := UpdateT{}
	err := json.Unmarshal(response, &update)
	if err != nil {
		return update, err
	}

	return update, nil
}

func sendMessage(chatId int, text string) (SendMessageResponseT, error) {
	url := BaseTelegramUrl + "/bot" + telegramToken + "/" + sendMessageUrl
	url = url + "?chat_id=" + strconv.Itoa(chatId) + "&text=" + text
	response := getResponse(url)

	sendMessage := SendMessageResponseT{}
	err := json.Unmarshal(response, &sendMessage)
	if err != nil {
		return sendMessage, err
	}
	return sendMessage, nil
}

func getResponse(url string) []byte {
	response := make([]byte, 0)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)

		return response
	}

	defer resp.Body.Close()

	for true {
		bs := make([]byte, 1024)
		n, err := resp.Body.Read(bs)
		response = append(response, bs[:n]...)

		if n == 0 || err != nil {
			break
		}
	}

	return response
}