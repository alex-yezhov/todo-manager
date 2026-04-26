package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := strings.TrimSpace(r.FormValue("now"))
	dateStr := strings.TrimSpace(r.FormValue("date"))
	repeat := strings.TrimSpace(r.FormValue("repeat"))

	var (
		now time.Time
		err error
	)

	if nowStr == "" {
		now = normalizeDate(time.Now())
	} else {
		now, err = parseDate(nowStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(next))
}

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	repeat = strings.TrimSpace(repeat)
	if repeat == "" {
		return "", errors.New("правило повторения не указано")
	}

	start, err := parseDate(dstart)
	if err != nil {
		return "", err
	}

	now = normalizeDate(now)

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("неверный формат правила")
	}

	switch parts[0] {
	case "d":
		return nextByDays(now, start, parts)
	case "y":
		return nextByYears(now, start, parts)
	case "w":
		return nextByWeek(now, start, parts)
	case "m":
		return nextByMonth(now, start, parts)
	default:
		return "", errors.New("неподдерживаемый формат правила")
	}
}

func nextByDays(now, start time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", errors.New("неверный формат правила d")
	}

	n, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", errors.New("неверный интервал дней")
	}
	if n < 1 || n > 400 {
		return "", errors.New("интервал дней должен быть от 1 до 400")
	}

	date := start
	for {
		date = date.AddDate(0, 0, n)
		if date.After(now) {
			return formatDate(date), nil
		}
	}
}

func nextByYears(now, start time.Time, parts []string) (string, error) {
	if len(parts) != 1 {
		return "", errors.New("неверный формат правила y")
	}

	date := start
	for {
		date = date.AddDate(1, 0, 0)
		if date.After(now) {
			return formatDate(date), nil
		}
	}
}

func nextByWeek(now, start time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", errors.New("неверный формат правила w")
	}

	weekdays, err := parseWeekdays(parts[1])
	if err != nil {
		return "", err
	}

	date := start.AddDate(0, 0, 1)

	for i := 0; i < 366*10; i++ {
		if date.After(now) && weekdays[weekdayNumber(date)] {
			return formatDate(date), nil
		}
		date = date.AddDate(0, 0, 1)
	}

	return "", errors.New("не удалось найти следующую дату")
}

func nextByMonth(now, start time.Time, parts []string) (string, error) {
	if len(parts) != 2 && len(parts) != 3 {
		return "", errors.New("неверный формат правила m")
	}

	days, err := parseMonthDays(parts[1])
	if err != nil {
		return "", err
	}

	months, err := parseMonths(parts[1:])
	if err != nil {
		return "", err
	}

	date := start.AddDate(0, 0, 1)

	for i := 0; i < 366*20; i++ {
		if date.After(now) && months[int(date.Month())] && matchesMonthDay(date, days) {
			return formatDate(date), nil
		}
		date = date.AddDate(0, 0, 1)
	}

	return "", errors.New("не удалось найти следующую дату")
}

func parseWeekdays(s string) ([8]bool, error) {
	var result [8]bool

	items := strings.Split(s, ",")
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			return result, errors.New("неверный список дней недели")
		}

		n, err := strconv.Atoi(item)
		if err != nil {
			return result, errors.New("неверный день недели")
		}
		if n < 1 || n > 7 {
			return result, errors.New("день недели должен быть от 1 до 7")
		}

		result[n] = true
	}

	return result, nil
}

func weekdayNumber(t time.Time) int {
	n := int(t.Weekday())
	if n == 0 {
		return 7
	}
	return n
}

func parseMonthDays(s string) ([]int, error) {
	items := strings.Split(s, ",")
	result := make([]int, 0, len(items))

	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			return nil, errors.New("неверный список дней месяца")
		}

		n, err := strconv.Atoi(item)
		if err != nil {
			return nil, errors.New("неверный день месяца")
		}

		if n == -1 || n == -2 || (n >= 1 && n <= 31) {
			result = append(result, n)
			continue
		}

		return nil, errors.New("день месяца должен быть 1..31, -1 или -2")
	}

	return result, nil
}

func parseMonths(parts []string) ([13]bool, error) {
	var result [13]bool

	if len(parts) == 1 {
		for i := 1; i <= 12; i++ {
			result[i] = true
		}
		return result, nil
	}

	if len(parts) != 2 {
		return result, errors.New("неверный список месяцев")
	}

	items := strings.Split(parts[1], ",")
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			return result, errors.New("неверный список месяцев")
		}

		n, err := strconv.Atoi(item)
		if err != nil {
			return result, errors.New("неверный месяц")
		}
		if n < 1 || n > 12 {
			return result, errors.New("месяц должен быть от 1 до 12")
		}

		result[n] = true
	}

	return result, nil
}

func matchesMonthDay(t time.Time, days []int) bool {
	day := t.Day()
	last := lastDayOfMonth(t)

	for _, d := range days {
		switch d {
		case -1:
			if day == last {
				return true
			}
		case -2:
			if day == last-1 {
				return true
			}
		default:
			if day == d {
				return true
			}
		}
	}

	return false
}

func lastDayOfMonth(t time.Time) int {
	y, m, _ := t.Date()
	return time.Date(y, m+1, 0, 0, 0, 0, 0, time.Local).Day()
}
