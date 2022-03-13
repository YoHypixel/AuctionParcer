package main

import (
	"fmt"
	"net/http"
	"sync"
)

type AuctionData struct {
	Success       bool      `json:"success"`
	Page          int       `json:"page"`
	TotalPages    int       `json:"totalPages"`
	TotalAuctions int       `json:"totalAuctions"`
	LastUpdated   int64     `json:"lastUpdated"`
	Auctions      []Auction `json:"auctions"`
}

type Auction struct {
	Uuid             string   `json:"uuid"`
	Auctioneer       string   `json:"auctioneer"`
	ProfileId        string   `json:"profile_id"`
	Coop             []string `json:"coop"`
	Start            int64    `json:"start"`
	End              int64    `json:"end"`
	ItemName         string   `json:"item_name"`
	ItemLore         string   `json:"item_lore"`
	Extra            string   `json:"extra"`
	Category         string   `json:"category"`
	Tier             string   `json:"tier"`
	StartingBid      int      `json:"starting_bid"`
	Claimed          bool     `json:"claimed"`
	HighestBidAmount int      `json:"highest_bid_amount"`
	Bin              bool     `json:"bin,omitempty"`
}

func NewClient() *http.Client {
	tr := &http.Transport{
		TLSHandshakeTimeout:   0,
		ResponseHeaderTimeout: 0,
		IdleConnTimeout:       0,
		MaxConnsPerHost:       0,
		MaxIdleConnsPerHost:   2000000,
	}

	c := http.Client{
		Transport: tr,
	}

	return &c
}

func BasicData() (*http.Client, AuctionData, error) {
	newClient := NewClient()
	auction, err := AuctionRequest(0, newClient)
	if err != nil {
		fmt.Println(err)
		return newClient, AuctionData{}, err
	}
	return newClient, auction, nil
}

type ItemData struct {
	Mutex          sync.Mutex
	ItemNames      []string
	Auctions       []Auction
	SortedAuctions map[string][]Auction
}

func (c *ItemData) AddData(waitGroup *sync.WaitGroup, page int, client *http.Client) (AuctionData, error) {
	current, err := AuctionRequest(page, client)
	defer waitGroup.Done()
	if err != nil {
		return AuctionData{}, err
	}

	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	for _, i := range current.Auctions {

		c.Auctions = append(c.Auctions, i)
	}
	return current, nil
}
