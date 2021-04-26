package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/kr/pretty"
	"github.com/rs/zerolog/log"
)

// https://maps.googleapis.com/maps/api/place/textsearch/json?key=AIzaSyC3MrrX8-db5JZW2_LwhwsjN1yXjdkj5YQ

type (
	queryPlaceResponse struct {
		Results []struct {
			PlaceID string `json:"place_id"`
		} `json:"results"`
	}

	placeAPI struct {
		key  string
		lang string
	}
)

const placeTextSearch = "https://maps.googleapis.com/maps/api/place/textsearch/json"

func (a *placeAPI) doGet(reqURL string, query url.Values) ([]byte, error) {
	if query == nil {
		query = make(url.Values)
	}

	query.Set("key", a.key)
	query.Set("language", a.lang)

	res, err := http.Get(reqURL + "?" + query.Encode())
	if err != nil {
		return nil, fmt.Errorf("unable to do request to Google: %w", err)
	}
	defer res.Body.Close()

	bb, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read request body from Google: %w", err)
	}

	return bb, nil
}

func main() {
	address := "Москва, Тверская, 14"

	api := &placeAPI{
		key:  "AIzaSyC3MrrX8-db5JZW2_LwhwsjN1yXjdkj5YQ",
		lang: "ru",
	}

	bb, err := api.doGet(placeTextSearch, url.Values{"input": []string{address}})
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка обращения в Google")
	}

	qr := &queryPlaceResponse{}

	err = json.Unmarshal(bb, qr)
	if err != nil {
		log.Fatal().Err(err).Msg("Невозможно демаршализовать JSON из ответа Google")
	}

	pretty.Println(qr)

	fmt.Fprintln(os.Stdout, string(bb))
}
