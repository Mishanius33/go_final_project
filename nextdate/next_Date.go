package nextdate

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	HoursPerDay = 24
	DateFormat  = "20060102"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	// Преобразование исходной даты
	t, err := time.Parse(DateFormat, date)
	if err != nil {
		return "", fmt.Errorf("Ошибка преобразования даты: %v", err)
	}

	// Отбрасываем время из now
	now = now.Truncate(HoursPerDay * time.Hour)

	// Проверка правила повторения
	switch {
	case repeat == "":
		return "", fmt.Errorf("Правило повторения не указано")

	case strings.HasPrefix(repeat, "d "):
		days, err := strconv.Atoi(strings.TrimPrefix(repeat, "d "))
		if err != nil || days < 1 || days > 400 {
			return "", fmt.Errorf("Неверный формат правила повторения: %s", repeat)
		}

		log.Printf("Начальная дата: %s, правило повторения: %d дней\n", t.Format(DateFormat), days)

		// Проверяю, совпадает ли исходная дата с текущей
		if t.Equal(now) {
			log.Printf("Дата совпадает с текущей: %s\n", t.Format(DateFormat))
			return now.AddDate(0, 0, days).Format(DateFormat), nil
		}

		// Добавляю дату перед началом цикла
		t = t.AddDate(0, 0, days)

		// Проверяю, стала ли дата больше текущей
		if t.After(now) {
			log.Printf("Дата сьала больше текущей: %s\n", t.Format(DateFormat))
			return t.Format(DateFormat), nil
		}

		// Если нет, начинаю цикл
		for {
			t = t.AddDate(0, 0, days)
			if t.After(now) {
				log.Printf("Дата после цикла: %s\n", t.Format(DateFormat))
				return t.Format(DateFormat), nil
			}
		}

	case repeat == "y":
		for {
			nextDate := t.AddDate(1, 0, 0)
			if nextDate.After(now) {
				return nextDate.Format(DateFormat), nil
			}
			t = nextDate
		}

	case strings.HasPrefix(repeat, "w "):
		weekdays := make(map[int]bool)
		for _, wd := range strings.Split(strings.TrimPrefix(repeat, "w "), ",") {
			wdNum, err := strconv.Atoi(wd)
			if err != nil || wdNum < 1 || wdNum > 7 {
				return "", fmt.Errorf("Неверный формат правила повторения: %s", repeat)
			}
			weekdays[wdNum] = true
		}
		for {
			t = t.AddDate(0, 0, 1)
			if weekdays[int(t.Weekday())+1] {
				return t.Format(DateFormat), nil
			}
		}

	default:
		return "", fmt.Errorf("Неверный формат правила повторения: %s", repeat)
	}
}
