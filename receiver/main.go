package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/coderero/toll_calculator/types"
	"github.com/gorilla/websocket"
)

type reciver struct {
	upgrader *websocket.Upgrader
	conn     *websocket.Conn
	kafka    DataProducer
}

func newReciver() (*reciver, error) {
	var (
		kafkaProducer DataProducer
		err           error
		kafkaTopic    = "obu_data"
	)

	if kafkaProducer, err = NewKafkaProducer(kafkaTopic); err != nil {
		return nil, err
	}

	kafkaProducer = &LoggingMiddleware{next: kafkaProducer}

	return &reciver{
		upgrader: &websocket.Upgrader{
			CheckOrigin:     func(r *http.Request) bool { return true },
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		kafka: kafkaProducer,
	}, nil

}

func (dr *reciver) produceData(data types.OBUData) error {
	return dr.kafka.Produce(data)
}

func (dr *reciver) wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := dr.upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	dr.conn = conn

	go dr.wsReadLoop()
}

func (dr *reciver) wsReadLoop() {
	for {
		var obuData types.OBUData
		if err := dr.conn.ReadJSON(&obuData); err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived, websocket.CloseAbnormalClosure) {
				dr.conn.Close()
				return
			}

			log.Fatal(err)

			var syntaxErr *json.SyntaxError
			if errors.As(err, &syntaxErr) {
				dr.conn.WriteJSON(map[string]string{"error": err.Error()})
			}
			continue
		}
		if err := dr.produceData(obuData); err != nil {
			log.Fatal(err)
		}

	}
}

func main() {
	dr, err := newReciver()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/ws", dr.wsHandler)
	if err := http.ListenAndServe(":8000", nil); err != nil {
		fmt.Println(err)
	}
}
