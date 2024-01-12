package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	nip05 := GetNip05FromUser()
	nip05info, err := GetPubkeyAndRelays(nip05)
	if err != nil {
		fmt.Println("Error getting relay information:", err)
		return
	}

	timer := time.NewTimer(10 * time.Second)
	defer timer.Stop()

	for _, url := range nip05info.Relays {
		fmt.Println("Connecting to relay:", url)

		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		defer conn.Close()

		if err != nil {
			fmt.Println("Error connecting to WebSocket:", err)
			return
		}

		reqId := "events_for_" + nip05
		req := []interface{}{
			"REQ",
			reqId,
			map[string]interface{}{
				"authors": []string{nip05info.Pubkey},
			},
		}

		if err := conn.WriteJSON(req); err != nil {
			fmt.Println("Error sending request:", err)
			return
		}

		for {
			exitReached := false
			select {
			case <-timer.C:
				fmt.Println("10 seconds passed without events, moving to next relay")
				exitReached = true
				break
			default:
				_, message, err := conn.ReadMessage()
				if err != nil {
					fmt.Println("Error reading message:", err)
					return
				}

				var msg []interface{}
				if err := json.Unmarshal(message, &msg); err != nil {
					fmt.Println("Error decoding message:", err)
					continue
				}

				if len(msg) < 2 {
					continue
				}

				if msg[0] == "EOSE" && msg[1] == reqId {
					fmt.Println("Received EOSE")
					exitReached = true
					break
				}

				if msg[0] == "EVENT" && len(msg) == 3 {
					eventData, err := json.Marshal(msg[2])
					if err != nil {
						fmt.Println("Error decoding event:", err)
						continue
					}

					nostrEvent := NostrEvent{}
					if err := json.Unmarshal(eventData, &nostrEvent); err != nil {
						fmt.Println("Error decoding event:", err)
						continue
					}
					SaveEventToArchive(nostrEvent)
				}
				if !timer.Stop() {
					<-timer.C
				}
				fmt.Println("Resetting timer")
				timer.Reset(10 * time.Second)
			}
			if exitReached {
				break
			}
		}
		conn.Close()
	}
}
