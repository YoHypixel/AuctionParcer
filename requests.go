package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
)

func AuctionRequest(page int, client *http.Client) (AuctionData, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.hypixel.net/skyblock/auctions", nil)
	if err != nil {
		log.Fatalln(err)
	}

	if client == nil {
		log.Fatalln("Found it")
	}

	params := req.URL.Query()
	params.Add("page", strconv.Itoa(page))
	req.URL.RawQuery = params.Encode()

	resp, err := client.Do(req)

	if err != nil {
		fmt.Printf("Error with request: %v\n", err)
		return AuctionData{}, errors.New("error doing request")

	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status)

		return AuctionData{}, errors.New("request is bad ")
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			log.Panicf("unable to close body: %v\n", err)
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return AuctionData{}, errors.New("unable to read request body")

	}

	var data AuctionData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return AuctionData{}, errors.New("unable to parse data")

	}

	return data, nil
}

// GetAllItemNames gets data from all pages
func GetAllItemNames() (*ItemData, AuctionData, error) {
	client, data, err := BasicData()

	if err != nil {
		return &ItemData{}, data, errors.New("unable to get client")
	}

	itemData := ItemData{}
	var wg sync.WaitGroup
	for i := 0; i < data.TotalPages; i++ {

		wg.Add(1)

		err = itemData.AddData(&wg, i, client)

		if err != nil {
			return &itemData, data, err
		}

	}
	wg.Wait()
	return &itemData, data, nil
}
