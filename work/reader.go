package work

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"sync/atomic"
)

var (
	count                = 0
	url                  = "https://geo-pg.captainjustin.space/transactions"
	skip                 = 0
	ops           uint64 = 0
	channels      []chan *Item
	lastChannelId = -1
)

func Start() {

	for i := 0; i < 5; i++ {
		c := createWorker()
		channels = append(channels, c)
	}
	fmt.Println("Channel count is:", len(channels))

	readAll()

}

func Block() {

}

func readAll() {

	filename := "items.data"

	f, err := os.Open(filename)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		os.Exit(1)
	}
	r := bufio.NewReader(f)
	strDate, e := Readln(r)
	strGeo, e := Readln(r)
	for e == nil {
		handle(strDate, strGeo)
		strDate, e = Readln(r)
		strGeo, e = Readln(r)
	}
}

func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

func handle(strDate string, strGeo string) {

	i := newItem(strDate, strGeo)

	lastChannelId++
	if lastChannelId >= len(channels) {
		lastChannelId = 0
	}

	channels[lastChannelId] <- i
	fmt.Println("Queued an item")
	atomic.AddUint64(&ops, 1)

}

type Item struct {
	Id     string `json:"id"`
	Lon    string `json:"lon"`
	Lat    string `json:"lat"`
	Tim    string `json:"tim"`
	Amount string `json:"amount"`
}

func newItem(strDate string, strLoc string) *Item {

	count++

	lo, la := explode(strLoc)

	i := Item{
		Id:     strconv.Itoa(count),
		Lon:    lo,
		Lat:    la,
		Tim:    strDate,
		Amount: "0.0"}

	return &i
}
