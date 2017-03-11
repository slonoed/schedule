package main

import (
	"html/template"
	"net/http"
	"time"
)

// User u
type User struct {
	Name string
	ID   int
}

// Day d
type Day struct {
	Time time.Time
	User *User
}

// Month m
type Month struct {
	Title string
	Days  []*Day
}

// Calendar all data
type Calendar struct {
	Users  []User
	Months []Month
}

var templates = template.Must(template.ParseGlob("templates/*"))

func isWorkday(day time.Time) bool {
	// add holidays
	return day.Weekday() != time.Sunday && day.Weekday() != time.Saturday
}

// prependEmpty добавляет пустые дни, если месяц начался не с ПН
func prependEmpty(month *Month) {
	first := month.Days[0].Time.Weekday()
	count := first - 1
	// Sunday
	if count == -1 {
		count = 6
	}
	for i := 0; i < int(count); i++ {
		month.Days = append([]*Day{nil}, month.Days...)
	}
}

func generateSchedule(names []string) Calendar {
	users := make([]User, len(names), len(names))
	for idx, name := range names {
		users[idx] = User{name, idx}
	}

	limit := 300
	months := []Month{}
	d := time.Now()
	currentMonth := Month{d.Month().String(), []*Day{&Day{d, &users[0]}}}
	i := 0
	userIdx := 1
	for i < limit {
		newD := d.AddDate(0, 0, 1)
		// Если новый месяц
		if newD.Month() != d.Month() {
			months = append(months, currentMonth)
			currentMonth = Month{newD.Month().String(), []*Day{}}
		}
		d = newD
		day := Day{d, nil}
		if isWorkday(d) {
			day.User = &users[(userIdx+1)%len(names)]
			userIdx++
		}
		currentMonth.Days = append(currentMonth.Days, &day)
		i++
	}

	months = append(months, currentMonth)
	for idx := range months {
		prependEmpty(&months[idx])
	}
	cal := Calendar{users, months}

	return cal
}

func indexPage(w http.ResponseWriter, req *http.Request) {
	templates.ExecuteTemplate(w, "indexPage", nil)
}

func resultPage(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	names := req.Form["username"]
	cal := generateSchedule(names)
	templates.ExecuteTemplate(w, "resultPage", cal)
}

func main() {
	http.HandleFunc("/", indexPage)
	http.HandleFunc("/result", resultPage)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.ListenAndServe(":6087", nil)
}
