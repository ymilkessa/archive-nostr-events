package main

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

func main() {
	nip05 := GetNip05FromUser()
	nip05info, err := GetPubkeyAndRelays(nip05)
	if err != nil {
		fmt.Println("Error getting relay information:", err)
		return
	}

	for _, url := range nip05info.Relays {

		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			fmt.Println("Error connecting to WebSocket:", err)
			return
		}
		defer conn.Close()

		const reqId string = "events_for_" + nip05
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
				break
			}

			if msg[0] == "EVENT" && len(msg) == 3 {

				if m, ok := msg[2].(map[string]interface{}); ok {
					event := NostrEvent{}
					event.Pubkey = m[EFPubkey].(string)
					event.Kind = int(m[EFKind].(float64))
					event.Id = m[EFId].(string)
					event.Content = m[EFContent].(string)
					event.CreatedAt = int64(m[EFCreatedAt].(float64))
					event.Sig = m[EFSig].(string)
					event.Tags = [][]string{}

					tagsInterface, ok := m[EFTags].([]interface{})
					if !ok {
						fmt.Println("Error decoding tags")
						continue
					} else {
						for _, tagInterface := range tagsInterface {
							tag, ok := tagInterface.([]interface{})
							if !ok {
								fmt.Println("Error decoding tag")
								continue
							} else {
								tagStr := []string{}
								for _, tagElem := range tag {
									tagStr = append(tagStr, tagElem.(string))
								}
								event.Tags = append(event.Tags, tagStr)
							}
						}
					}

					SaveEventToArchive(event)
				}
			}
		}
		conn.Close()
	}
}
