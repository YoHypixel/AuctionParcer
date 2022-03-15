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
		fmt.Printf("error with new http request %v\n", err)
	}

	req.Header.Set("user-agent", "auction parser golang")
	if page != 0 {
		req.URL.RawQuery = "page=" + strconv.Itoa(page)
	}
	resp, err := client.Do(req)

	if err != nil {

		fmt.Printf("Error with request: %v\n", err)
		if resp == nil {

			fmt.Printf("Nothing returned\n")

		} else {
			fmt.Printf("This is what was returned %v\n", resp.Status)
		}

		return AuctionData{}, errors.New("error doing request")
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Status)

		return AuctionData{}, errors.New("request is bad ")
	}
	fmt.Printf("Request succeeded, page: %v \n", req.URL)

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
	// client, data, err := BasicData()

	// if err != nil {
	//	return &ItemData{}, data, err
	//}

	itemData := ItemData{}
	var wg sync.WaitGroup

	testClient := NewClient()

	basicData, err := AuctionRequest(0, testClient)

	if err != nil {
		return &itemData, AuctionData{}, err
	}

	for i := 1; i < basicData.TotalPages; i++ {

		wg.Add(1)

		checkData, err := itemData.AddData(&wg, i, testClient)

		if err != nil {
			return &itemData, basicData, err
		}

		if checkData.TotalPages < basicData.TotalPages {
			basicData.TotalPages = checkData.TotalPages

		} else if checkData.TotalPages > basicData.TotalPages {
			basicData.TotalPages = checkData.TotalPages
		}
		// time.Sleep(50 * time.Millisecond)

	}
	wg.Wait()
	return &itemData, basicData, nil
}
