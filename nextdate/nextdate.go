package nextdate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {
	// Преобразование исходной даты
	t, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("Ошибка преобразования даты: %v", err)
	}

	// Отбрасываем время из now
	now = now.Truncate(24 * time.Hour)

	// Проверка правила повторения
	switch {
	case repeat == "":
		return "", fmt.Errorf("Правило повторения не указано")

	case strings.HasPrefix(repeat, "d "):
		days, err := strconv.Atoi(strings.TrimPrefix(repeat, "d "))
		if err != nil || days < 1 || days > 400 {
			return "", fmt.Errorf("Неверный формат правила повторения: %s", repeat)
		}

		// Проверяю, совпадает ли исходная дата с текущей
		if t.Equal(now) {
			return t.Format("20060102"), nil
		}

		// Добавляю дату перед началом цикла
		t = t.AddDate(0, 0, days)

		// Проверяю, стала ли дата больше текущей
		if t.After(now) {
			return t.Format("20060102"), nil
		}

		// Если нет, начинаю цикл
		for {
			t = t.AddDate(0, 0, days)
			if t.After(now) {
				return t.Format("20060102"), nil
			}
		}

	case repeat == "y":
		for {
			nextDate := t.AddDate(1, 0, 0)
			if nextDate.After(now) {
				return nextDate.Format("20060102"), nil
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
				return t.Format("20060102"), nil
			}
		}

	case strings.HasPrefix(repeat, "m "):
		parts := strings.Split(strings.TrimPrefix(repeat, "m "), " ")
		if len(parts) < 1 || len(parts) > 2 {
			return "", fmt.Errorf("Неверный формат правила повторения: %s", repeat)
		}
		days := make([]int, 0, len(parts[0]))
		for _, d := range strings.Split(parts[0], ",") {
			var day int
			if d == "-1" {
				day = -1
			} else if d == "-2" {
				day = -2
			} else {
				day, err = strconv.Atoi(d)
				if err != nil || (day != -1 && day != -2 && (day < 1 || day > 31)) {
					return "", fmt.Errorf("Неверный формат правила повторения: %s", repeat)
				}
			}
			days = append(days, day)
		}
		months := make([]int, 0, 12)
		if len(parts) == 2 {
			for _, m := range strings.Split(parts[1], ",") {
				month, err := strconv.Atoi(m)
				if err != nil || month < 1 || month > 12 {
					return "", fmt.Errorf("Неверный формат правила повторения: %s", repeat)
				}
				months = append(months, month)
			}
		} else {
			for i := 1; i <= 12; i++ {
				months = append(months, i)
			}
		}
		for {
			for _, m := range months {
				for _, d := range days {
					nextDate := time.Date(t.Year(), time.Month(m), d, 0, 0, 0, 0, t.Location())
					if d < 0 {
						nextDate = nextDate.AddDate(0, 1, 0).AddDate(0, 0, d)
					}
					if nextDate.After(now) {
						return nextDate.Format("20060102"), nil
					}
				}
			}
			t = t.AddDate(0, 1, 0)
		}
	default:
		return "", fmt.Errorf("Неверный формат правила повторения: %s", repeat)
	}
}
