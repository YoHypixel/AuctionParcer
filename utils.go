package main

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
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
		MaxIdleConnsPerHost: 2048,
		TLSHandshakeTimeout: 0 * time.Second,
	}

	c := http.Client{
		Transport: tr,
	}

	return &c
}

func BasicData() (*http.Client, AuctionData, error) {
	newClient := NewClient()

	auction, err := AuctionRequest(0)
	if err != nil {
		fmt.Println(err)
		return nil, AuctionData{}, errors.New("basic data")
	}
	return newClient, auction, nil
}

type ItemData struct {
	Mutex          sync.Mutex
	ItemNames      []string
	Auctions       []Auction
	SortedAuctions map[string][]Auction
}

func (c *ItemData) AddData(wg *sync.WaitGroup, page int) error {
	defer wg.Done()
	defer c.Mutex.Unlock()
	current, err := AuctionRequest(page)

	if err != nil {
		return err
	}

	c.Mutex.Lock()
	for _, i := range current.Auctions {

		if i.Uuid == "" {
			return errors.New("add data has failed")
		}

		c.Auctions = append(c.Auctions, i)
	}
	return nil
}
