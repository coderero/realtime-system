package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/coderero/toll_calculator/types"
	"github.com/gorilla/websocket"
)

const (
	sendInterval = 1
	wsURL        = "ws://localhost:8000/ws"
)

func generateCords() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()
	return n + f
}

func generateLocation() (float64, float64) {
	return generateCords(), generateCords()
}

func generateOBUIDS(n int) []int {
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p = append(p, i)
	}
	return p
}

func main() {
	obuIDs := generateOBUIDS(20)
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer conn.Close()

	for {
		for i := 0; i < len(obuIDs); i++ {
			lat, long := generateLocation()
			obuData := types.OBUData{
				OBUID: obuIDs[i],
				Lat:   lat,
				Long:  long,
			}
			if err := conn.WriteJSON(obuData); err != nil {
				log.Println(err)
				return
			}
			fmt.Printf("Sent: %v\n", obuData)
		}
		time.Sleep(sendInterval * time.Second)
	}
}
