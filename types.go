package main

type NostrEvent struct {
	Pubkey    string     `json:"pubkey"`
	Kind      int        `json:"kind"`
	Id        string     `json:"id"`
	Tags      [][]string `json:"tags"`
	Content   string     `json:"content"`
	CreatedAt int64      `json:"created_at"`
	Sig       string     `json:"sig"`
}
