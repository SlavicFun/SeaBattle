package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/time/rate"
)

type Pixel struct {
	Emoji  string
	Option string
	ID     string
	Type   string
	Ship   bool
	Fire   bool
	Death  bool
}

type Ship struct {
	Length int
	Coords []string
	Death  bool
}

func isAreaAroundFree(pixels []Pixel, newRow, newCol int) bool {
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if newRow+i >= 0 && newRow+i < 8 && newCol+j >= 0 && newCol+j < 8 && pixels[(newRow+i)*8+(newCol+j)].Ship {
				return false
			}
		}
	}
	return true
}

func RandomCoords(ships []Ship, pixels []Pixel) []Ship {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < len(ships); i++ {
		for {
			row := r.Intn(8)
			col := r.Intn(8)
			direction := r.Intn(2)

			valid := true
			coords := []string{}

			for j := 0; j < ships[i].Length; j++ {
				var newRow, newCol int

				if direction == 0 {
					newRow = row
					newCol = col + j
				} else {
					newRow = row + j
					newCol = col
				}

				if newRow < 0 || newRow >= 8 || newCol < 0 || newCol >= 8 || !isAreaAroundFree(pixels, newRow, newCol) {
					valid = false
					break
				}

				coord := fmt.Sprintf("%c%d", 'A'+newRow, newCol+1)
				coords = append(coords, coord)
			}

			if valid {
				ships[i].Coords = coords
				for _, coord := range coords {
					pixels[int((coord[0]-'A')*8)+(int(coord[1])-'1')].Ship = true
					pixels[int((coord[0]-'A')*8)+(int(coord[1])-'1')].Emoji = strconv.Itoa(ships[i].Length) + "-" + strconv.Itoa(i)

				}
				break
			}
		}
	}

	return ships
}

func GetPixels2() (Pixels []Pixel) {
	ships := []Ship{
		{Length: 1, Coords: []string{}},
		{Length: 1, Coords: []string{}},
		{Length: 1, Coords: []string{}},
		{Length: 2, Coords: []string{}},
		{Length: 2, Coords: []string{}},
		{Length: 3, Coords: []string{}},
		{Length: 4, Coords: []string{}},
	}

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			coord := fmt.Sprintf("%c%d", 'A'+i, j+1)
			Pixels = append(Pixels, Pixel{ID: coord, Fire: false, Ship: false, Emoji: " ", Option: fmt.Sprintf("option_%d", i*8+j+1)})
		}
	}

	RandomCoords(ships, Pixels)

	return Pixels
}

func CheckShipDeath(Pixels []Pixel, Emoji string) []Pixel {
	for i := range Pixels {
		if Pixels[i].Emoji == Emoji {
			Pixels[i].Death = true
		}
	}
	Pixels = FindEmptyAdjacentCells(Pixels, Emoji)
	return Pixels
}

func FindEmptyAdjacentCells(Pixels []Pixel, emoji string) []Pixel {

	directions := []struct{ dx, dy int }{
		{-1, -1}, {-1, 0}, {-1, 1},
		{0, -1}, {0, 1},
		{1, -1}, {1, 0}, {1, 1},
	}

	// Находим все клетки с заданным Emoji и их координаты
	for i, pixel := range Pixels {
		if pixel.Emoji == emoji {
			x, y := i/8, i%8

			for _, d := range directions {
				newX, newY := x+d.dx, y+d.dy
				if newX >= 0 && newX < 8 && newY >= 0 && newY < 8 {
					index := newX*8 + newY
					if !Pixels[index].Ship {
						Pixels[index].Fire = true
					}
				}
			}
		}
	}

	return Pixels
}

func PrintMap(Pixels []Pixel) tgbotapi.InlineKeyboardMarkup {
	var Row1Pixels []Pixel
	var Row2Pixels []Pixel
	var Row3Pixels []Pixel
	var Row4Pixels []Pixel
	var Row5Pixels []Pixel
	var Row6Pixels []Pixel
	var Row7Pixels []Pixel
	var Row8Pixels []Pixel

	var row1 []tgbotapi.InlineKeyboardButton
	var row2 []tgbotapi.InlineKeyboardButton
	var row3 []tgbotapi.InlineKeyboardButton
	var row4 []tgbotapi.InlineKeyboardButton
	var row5 []tgbotapi.InlineKeyboardButton
	var row6 []tgbotapi.InlineKeyboardButton
	var row7 []tgbotapi.InlineKeyboardButton
	var row8 []tgbotapi.InlineKeyboardButton

	// Создаем первый ряд кнопок
	row1 = []tgbotapi.InlineKeyboardButton{}
	Row1Pixels = Pixels[0:8]
	for i := 0; i < len(Row1Pixels); i++ {
		row1 = append(row1, tgbotapi.NewInlineKeyboardButtonData(GetRowEmojy(Row1Pixels[i]), Row1Pixels[i].Option))
	}
	// Создаем второй ряд кнопок
	row2 = []tgbotapi.InlineKeyboardButton{}
	Row2Pixels = Pixels[8:16]
	for i := 0; i < len(Row2Pixels); i++ {
		row2 = append(row2, tgbotapi.NewInlineKeyboardButtonData(GetRowEmojy(Row2Pixels[i]), Row2Pixels[i].Option))
	}
	// Создаем третий ряд кнопок
	row3 = []tgbotapi.InlineKeyboardButton{}
	Row3Pixels = Pixels[16:24]
	for i := 0; i < len(Row3Pixels); i++ {
		row3 = append(row3, tgbotapi.NewInlineKeyboardButtonData(GetRowEmojy(Row3Pixels[i]), Row3Pixels[i].Option))
	}
	// Создаем четвертый ряд кнопок
	row4 = []tgbotapi.InlineKeyboardButton{}
	Row4Pixels = Pixels[24:32]
	for i := 0; i < len(Row4Pixels); i++ {
		row4 = append(row4, tgbotapi.NewInlineKeyboardButtonData(GetRowEmojy(Row4Pixels[i]), Row4Pixels[i].Option))
	}
	// Создаем пятый ряд кнопок
	row5 = []tgbotapi.InlineKeyboardButton{}
	Row5Pixels = Pixels[32:40]
	for i := 0; i < len(Row5Pixels); i++ {
		row5 = append(row5, tgbotapi.NewInlineKeyboardButtonData(GetRowEmojy(Row5Pixels[i]), Row5Pixels[i].Option))
	}
	// Создаем шестой ряд кнопок
	row6 = []tgbotapi.InlineKeyboardButton{}
	Row6Pixels = Pixels[40:48]
	for i := 0; i < len(Row6Pixels); i++ {
		row6 = append(row6, tgbotapi.NewInlineKeyboardButtonData(GetRowEmojy(Row6Pixels[i]), Row6Pixels[i].Option))
	}
	// Создаем седьмой ряд кнопок
	row7 = []tgbotapi.InlineKeyboardButton{}
	Row7Pixels = Pixels[48:56]
	for i := 0; i < len(Row7Pixels); i++ {
		row7 = append(row7, tgbotapi.NewInlineKeyboardButtonData(GetRowEmojy(Row7Pixels[i]), Row7Pixels[i].Option))
	}
	// Создаем восьмой ряд кнопок
	row8 = []tgbotapi.InlineKeyboardButton{}
	Row8Pixels = Pixels[56:64]
	for i := 0; i < len(Row8Pixels); i++ {

		row8 = append(row8, tgbotapi.NewInlineKeyboardButtonData(GetRowEmojy(Row8Pixels[i]), Row8Pixels[i].Option))
	}

	// Создаем клавиатуру
	return tgbotapi.NewInlineKeyboardMarkup(row1, row2, row3, row4, row5, row6, row7, row8)
}

func StartBattle() {

	// Устанавливаем токен для бота
	bot, err := tgbotapi.NewBotAPI("5645589091:AAHd-8GFqt0tFe50We2gVBw_8VNrBOvI6r4")
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 34

	updates, _ := bot.GetUpdatesChan(u)

	// устанавливаем максимальную скорость отправки запросов в 5 запросов в секунду
	limiter := rate.NewLimiter(5, 1)

	Pixels := GetPixels2()

	var keyboard tgbotapi.InlineKeyboardMarkup

	for update := range updates {
		if !limiter.Allow() {
			log.Println("Rate Limit: ", !limiter.Allow())
			continue
		}
		if update.CallbackQuery != nil {
			callback := update.CallbackQuery
			status := true
			game := false
			for i := range Pixels {
				log.Println("lol")
				if Pixels[i].Option == callback.Data {
					log.Println("lol2")
					if Pixels[i].Ship && Pixels[i].Fire {
						status = false
					} else if Pixels[i].Ship && !Pixels[i].Fire {
						Pixels[i].Fire = true
						UnfireFire := 0
						for j := range Pixels {
							if Pixels[j].Emoji == Pixels[i].Emoji {
								if !Pixels[j].Fire {
									UnfireFire++
								}
							}
						}
						if UnfireFire == 0 {
							Pixels = CheckShipDeath(Pixels, Pixels[i].Emoji)
						}
					} else if !Pixels[i].Ship && Pixels[i].Fire {
						status = false
					} else if !Pixels[i].Ship && !Pixels[i].Fire {
						Pixels[i].Fire = true
					} else {
						status = false
					}
				}
				if Pixels[i].Ship && !Pixels[i].Fire {
					log.Println("lol3")
					game = true
				}
			}
			if !game {
				log.Println("lol4")
				Pixels = GetPixels2()
			}
			if status {
				msg := tgbotapi.NewEditMessageText(callback.Message.Chat.ID, callback.Message.MessageID, "Ваше поле:")

				keyboard = PrintMap(Pixels)
				msg.ReplyMarkup = &keyboard

				// Отправляем сообщение
				_, err := bot.Send(msg)
				if err != nil {
					if err.Error() == "Bad Request: message is not modified" {
						log.Println("Message is not modified")
					} else {
						log.Panic(err)
					}
				}
			}

		} else if update.Message != nil {
			log.Println("kek")
			NewMsg := tgbotapi.NewMessage(update.Message.Chat.ID, "Поле соперника:")

			// Создаем инлайн клавиатуру и добавляем в нее ряды кнопок
			keyboard = PrintMap(Pixels)

			// Устанавливаем созданную клавиатуру в сообщении
			NewMsg.ReplyMarkup = keyboard

			_, err := bot.Send(NewMsg)
			if err != nil {
				if err.Error() == "Bad Request: message is not modified" {
					log.Println("Message is not modified")
				} else {
					log.Panic(err)
				}
			}
		} else {
			continue
		}
	}
}

func GetRowEmojy(Pixel Pixel) string {
	if Pixel.Death {
		return "💀"
	} else if Pixel.Ship {
		if Pixel.Fire {
			return "🔥"
		} else {
			return " "
			// return "⛴"
		}
	} else {
		if Pixel.Fire {
			return "🗯"
		} else {
			return " "
		}
	}

}
