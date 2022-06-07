package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const timeoutSec = 5

type blockResponse []struct {
	Hash       string `json:"hash"`
	Height     int    `json:"height"`
	Time       int64  `json:"time"`
	BlockIndex int    `json:"block_index"`
}

func getBlockData(tim time.Time) (out resp, err error) {
	var (
		data blockResponse
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*timeoutSec))
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://blockchain.info/blocks/%d?format=json", tim.Add(time.Hour).UnixMilli()), nil)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return
	}

	body, _ := ioutil.ReadAll(res.Body)
	defer func() {
		_ = res.Body.Close()
	}()
	if err = json.Unmarshal(body, &data); err != nil {
		return
	}
	for i := 1; i < len(data); i++ {
		if data[i].Time < tim.Unix() {
			out.Top = data[i-1].BlockIndex
			out.Bottom = data[i].BlockIndex
			parc := float64((tim.Unix() - data[i].Time)) / float64((data[i-1].Time - data[i].Time))
			out.Middle = float64(data[i].BlockIndex) + parc
			return
		}
	}
	return resp{}, fmt.Errorf("didn't find block data for this time")
}
