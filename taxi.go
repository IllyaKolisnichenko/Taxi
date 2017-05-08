package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

var muxLock sync.Mutex

const sizeSymbol int = 25
const sizeRequest int = 50

//структура информации администратору
type application struct {
	title string
	views int
}

var numberAdminRequest int = 0
var firstPart, secondPart int

//информация по заявкам администратору
var ArrRequest map[int]application = make(map[int]application)

//активные заявки
var request map[int]string = make(map[int]string)

//map с символами,для создания заявок
var symbol map[int]string = map[int]string{0: "a", 1: "b", 2: "c", 3: "d", 4: "e",
	5: "f", 6: "j", 7: "h", 8: "i", 9: "g",
	10: "k", 11: "l", 12: "m", 13: "n", 14: "o",
	15: "p", 16: "q", 17: "r", 18: "s", 19: "t",
	20: "u", 21: "v", 22: "w", 23: "x", 24: "y",
	25: "z"}

func sort() { //подсчет количества выводов заявок таксисту
	for i := 0; i < len(ArrRequest); i++ {
		for j := i + 1; j < len(ArrRequest); j++ {
			if ArrRequest[i].title == ArrRequest[j].title {
				muxLock.Lock()
				ArrRequest[j] = application{title: ArrRequest[i].title, views: ArrRequest[i].views + 1}
				ArrRequest[i] = application{views: 0}
				muxLock.Unlock()
			}
		}
	}
}

func mapFilling() { //функция начального заполнения мапа заявками
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < sizeRequest; i++ {
		muxLock.Lock()
		firstPart = rand.Intn(sizeSymbol)
		secondPart = rand.Intn(sizeSymbol)
		request[i] = symbol[firstPart] + symbol[secondPart]
		muxLock.Unlock()
	}
}

func replacement() { //функция замещения заявок раз в 200мс
	for {
		muxLock.Lock()
		var number int = rand.Intn(sizeRequest)
		firstPart = rand.Intn(sizeSymbol)
		secondPart = rand.Intn(sizeSymbol)
		var x string = symbol[firstPart] + symbol[secondPart]
		request[number] = x
		muxLock.Unlock()
		time.Sleep(time.Millisecond * 200)
	}
}
func main() {
	mapFilling()
	go replacement()
	router := mux.NewRouter()
	router.HandleFunc("/{request}", cabbie).Methods("GET")
	router.HandleFunc("/request/{admin}", admin).Methods("GET")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func cabbie(w http.ResponseWriter, r *http.Request) { //ф-ция-обработчик запросов от таксистов
	vars := mux.Vars(r)
	interrogator := vars["request"]
	if interrogator == "request" {
		muxLock.Lock()
		x := rand.Intn(sizeRequest)
		fmt.Fprintln(w, "Заказ:", request[x])
		ArrRequest[numberAdminRequest] = application{title: request[x], views: ArrRequest[numberAdminRequest].views + 1}
		numberAdminRequest++
		muxLock.Unlock()
	}
}
func admin(w http.ResponseWriter, r *http.Request) { //функция-обработчик запросов от администратора
	vars := mux.Vars(r)
	interrogator := vars["admin"]
	sort()
	if interrogator == "admin" {
		for i := 0; i < len(ArrRequest); i++ {
			if ArrRequest[i].views > 0 {
				fmt.Fprintln(w, "Заказ:", ArrRequest[i].title, "Количество показов:", ArrRequest[i].views)
			}
		}
	}
}
