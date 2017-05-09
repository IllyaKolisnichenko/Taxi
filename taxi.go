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

//структура информации администратору
type application struct {
	title string
	views int
}

var (
	muxLock               sync.Mutex
	numberAdminRequest    int = 0
	firstPart, secondPart int
	//информация по заявкам администратору
	ArrRequest map[int]application = make(map[int]application)

	//активные заявки
	request map[int]string = make(map[int]string)
	//Arr с символами,для создания заявок
	symbol [sizeSymbol + 1]string = [sizeSymbol + 1]string{"a", "b", "c", "d", "e",
						               "f", "j", "h", "i", "g",
		                                               "k", "l", "m", "n", "o",
		                                               "p", "q", "r", "s", "t",
		                                               "u", "v", "w", "x", "y",
		                                               "z"}
)

const (
	usrAdmin    string = "admin"
	usrCabbie   string = "request"
	sizeSymbol  int    = 25
	sizeRequest int    = 50
)

func sort() { //подсчет количества выводов заявок таксисту
	for i := 0; i < len(ArrRequest); i++ {
		for j := i + 1; j < len(ArrRequest); j++ {
			if ArrRequest[i].title == ArrRequest[j].title {
				ArrRequest[j] = application{title: ArrRequest[i].title, views: ArrRequest[i].views + 1}
				ArrRequest[i] = application{views: 0}
			}
		}
	}
}

func mapFilling() { //функция начального заполнения мапа заявками
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < sizeRequest; i++ {
		firstPart = rand.Intn(sizeSymbol)
		secondPart = rand.Intn(sizeSymbol)
		muxLock.Lock()
		request[i] = symbol[firstPart] + symbol[secondPart]
		muxLock.Unlock()
	}
}

func replacement() { //функция замещения заявок раз в 200мс
	for {
		var number int = rand.Intn(sizeRequest)
		firstPart = rand.Intn(sizeSymbol)
		secondPart = rand.Intn(sizeSymbol)
		muxLock.Lock()
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
	interrogator := vars[usrCabbie]
	if interrogator == usrCabbie {
		x := rand.Intn(sizeRequest)
		muxLock.Lock()
		fmt.Fprintln(w, "Заказ:", request[x])
		muxLock.Unlock()
		muxLock.Lock()
		ArrRequest[numberAdminRequest] = application{title: request[x], views: ArrRequest[numberAdminRequest].views + 1}
		numberAdminRequest++
		muxLock.Unlock()
	}
}
func admin(w http.ResponseWriter, r *http.Request) { //функция-обработчик запросов от администратора
	vars := mux.Vars(r)
	interrogator := vars[usrAdmin]
	sort()
	if interrogator == usrAdmin {
		for i := 0; i < len(ArrRequest); i++ {
			if ArrRequest[i].views > 0 {
				fmt.Fprintln(w, "Заказ:", ArrRequest[i].title, "Количество показов:", ArrRequest[i].views)
			}
		}
	}
}
