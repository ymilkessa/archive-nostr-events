package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Nip05Response struct {
	Names  map[string]string   `json:"names"`
	Relays map[string][]string `json:"relays"`
}

type Nip05Info struct {
	Pubkey string
	Relays []string
}

const (
	nip05NamesField  string = "names"
	nip05RelaysField string = "relays"
)

// A function for getting a user's pubkey and their relays from a given nip05 identifier.
// See https://github.com/nostr-protocol/nips/blob/master/05.md on how this works.
// Returns a pointer to a Nip05Info struct instance:
//
//	type Nip05Info struct {
//		Pubkey string
//		Relays []string
//	}
func GetPubkeyAndRelays(nip05 string) (*Nip05Info, error) {
	at_index := strings.Index(nip05, "@")
	local_part := nip05[:at_index]
	domain := nip05[at_index+1:]

	request_url := fmt.Sprintf("https://%s/.well-known/nostr.json?name=%s", domain, local_part)
	resp, err := http.Get(request_url)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var data Nip05Response
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, err
	}
	pubkey := data.Names[local_part]
	relays := data.Relays[pubkey]
	nip05_info := Nip05Info{pubkey, relays}
	return &nip05_info, nil
}
