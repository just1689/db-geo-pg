package work

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync/atomic"
)

func createWorker() chan *Item {

	c := make(chan *Item, 5)

	go func() {
		for {
			i := <-c
			b := marshal(i)

			post(b)
			atomic.AddUint64(&pCount, 1)
		}

	}()

	fmt.Println("New worker listening!")
	return c
}

func marshal(i *Item) []byte {
	b, _ := json.Marshal(i)
	return b
}

func explode(strLonLat string) (lon string, lat string) {
	s := strings.Split(strLonLat, " ")
	return s[0], s[1]
}

func post(b []byte) {

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	_, _ = ioutil.ReadAll(resp.Body)
}
