package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"

	"github.com/kr/pretty"
	"github.com/rs/zerolog/log"
)

// https://maps.googleapis.com/maps/api/place/textsearch/json?key=AIzaSyC3MrrX8-db5JZW2_LwhwsjN1yXjdkj5YQ

type (
	queryPlaceResp struct {
		Results []struct {
			PlaceID string `json:"place_id"`
		} `json:"results"`
	}

	placeAPI struct {
		key  string
		lang string
	}

	photo struct {
		Ref string `json:"photo_reference"`
	}

	queryDetailsResp struct {
		Result struct {
			Photos []photo `json:"photos"`
		} `json:"result"`
	}
)

const placeTextSearch = "https://maps.googleapis.com/maps/api/place/textsearch/json"

func (api *placeAPI) doGet(reqURL string, query url.Values) ([]byte, error) {
	if query == nil {
		query = make(url.Values)
	}

	query.Set("key", api.key)
	query.Set("language", api.lang)

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

func (api *placeAPI) getPhotoUrl(photoRef string, maxwidth int) string {
	query := make(url.Values)
	query.Set("key", api.key)
	query.Set("photoreference", photoRef)
	query.Set("maxwidth", strconv.Itoa(maxwidth))
	return "https://maps.googleapis.com/maps/api/place/photo?" + query.Encode()
}

func main() {
	address := "Земля Франца Иосифа"

	api := &placeAPI{
		key:  "AIzaSyC3MrrX8-db5JZW2_LwhwsjN1yXjdkj5YQ",
		lang: "ru",
	}

	bb, err := api.doGet(placeTextSearch, url.Values{"input": []string{address}})
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка обращения в Google")
	}

	qrPl := &queryPlaceResp{}

	err = json.Unmarshal(bb, qrPl)
	if err != nil {
		log.Fatal().Err(err).Msg("Невозможно демаршализовать JSON из ответа Google")
	}

	fmt.Fprintln(os.Stdout, string(bb))
	urls := []string{}
	for _, r := range qrPl.Results {
		uu, err := api.getPhotoURLs(r.PlaceID)
		if err != nil {
			log.Fatal().Err(err).Msg("Невозможно получить фото по месту из Google")
		}
		urls = append(urls, uu...)
	}

	if len(urls) == 0 {
		log.Fatal().Err(err).Msg("Нет фото об этом месте")
	}

	comm := exec.Command("feh", urls...)
	err = comm.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("Ошибка при запуске feh")
	}
}

func (api *placeAPI) getPhotoURLs(placeID string) ([]string, error) {
	bb, err := api.doGet("https://maps.googleapis.com/maps/api/place/details/json",
		url.Values{
			"place_id": []string{placeID},
			"fields":   []string{"photo"},
		},
	)
	if err != nil {
		return nil, err
	}

	qrDt := &queryDetailsResp{}
	err = json.Unmarshal(bb, qrDt)
	if err != nil {
		return nil, err
	}

	fmt.Fprintln(os.Stdout, string(bb))
	pretty.Println(qrDt)

	urls := make([]string, len(qrDt.Result.Photos))
	for i := range qrDt.Result.Photos {
		urls[i] = api.getPhotoUrl(qrDt.Result.Photos[i].Ref, 3200)
	}
	return urls, nil
}
