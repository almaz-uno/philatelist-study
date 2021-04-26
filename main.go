package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/rs/zerolog/log"
)

// https://maps.googleapis.com/maps/api/place/textsearch/json?key=AIzaSyC3MrrX8-db5JZW2_LwhwsjN1yXjdkj5YQ

const placeTextSearch = "https://maps.googleapis.com/maps/api/place/textsearch/json"

func main() {
	uq := make(url.Values)

	uq.Set("key", "AIzaSyC3MrrX8-db5JZW2_LwhwsjN1yXjdkj5YQ")
	uq.Set("input", "Москва, Тверская, 14")

	reqURL := placeTextSearch + "?" + uq.Encode()

	log.Info().Msg(reqURL)

	res, err := http.Get(reqURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Невозможно получить места")
	}
	defer res.Body.Close()

	bb, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal().Err(err).Msg("Невозможно считать ответ из ответа")
	}

	fmt.Fprintln(os.Stdout, string(bb))
}
