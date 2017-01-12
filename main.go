package main

import (
	"bufio"
	"encoding/json"
	"github.com/centrifugal/centrifuge-go"
	"github.com/centrifugal/centrifugo/libcentrifugo/auth"
	"github.com/ivpusic/grpool"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const nbClient = 1000
const nbMessage = 100
const totalMsg = nbClient * nbMessage

func newConnection(n int) centrifuge.Centrifuge {

	secret := "admin"
	user := strconv.Itoa(n)
	timestamp := centrifuge.Timestamp()
	info := ""

	token := auth.GenerateClientToken(secret, user, timestamp, info)

	creds := &centrifuge.Credentials{
		User:      user,
		Timestamp: timestamp,
		Info:      info,
		Token:     token,
	}

	wsURL := "ws://ws.server.io/connection/websocket"
	c := centrifuge.NewCentrifuge(wsURL, creds, nil, centrifuge.DefaultConfig)

	err := c.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	return c
}

func main() {

	pool := grpool.NewPool(100, 50)

	defer pool.Release()
	pool.WaitCount(nbClient)

	started := time.Now()
	count := 0

	for i := 0; i < nbClient; i++ {

		pool.JobQueue <- func() {

			defer pool.JobDone()

			stri := strconv.Itoa(i)

			c := newConnection(i)
			sub, _ := c.Subscribe("test", nil)

			log.Printf("Client %d is connected", i)

			for j := 0; j < nbMessage; j++ {

				r := rand.Intn(1000)
				time.Sleep(time.Duration(r) * time.Millisecond)

				strj := strconv.Itoa(j)

				data := map[string]string{"message": stri + " : " + strj + "</br>"}
				dataBytes, _ := json.Marshal(data)

				sub.Publish(dataBytes)
				count += 1

			}
		}

		time.Sleep(70 * time.Millisecond)
	}

	log.Print("... press key for done ...")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

	elapsed := time.Since(started)
	log.Printf("Total clients %d", nbClient)
	log.Printf("Total messages %d", totalMsg)
	log.Printf("Real Total messages %d", count)
	log.Printf("Elapsed %s", elapsed)
	log.Printf("Msg/sec %d", (1000*totalMsg)/int(elapsed.Nanoseconds()/1000000))
}
