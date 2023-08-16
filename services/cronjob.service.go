package services

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const BASE_URL = "https://youtube.googleapis.com/youtube/v3/search?part=snippet&order=date&publishedAfter=2023-08-01T00%3A00%3A00Z&q=minecraft&type=video&key="

func jsonToMap(jsonStr string) map[string]interface{} {
	result := make(map[string]interface{})
	json.Unmarshal([]byte(jsonStr), &result)
	return result
}
func YouTubeCronJob() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file", err)
		return
	}
	key := os.Getenv("GCP_APIKEY")
	resp, err := http.Get(BASE_URL + key)
	if err != nil {
		log.Println("Error fetching from youtube", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	sb := string(body)
	res := jsonToMap(sb)
	data := res["items"]
	var values []map[string]interface{}
	for _, v := range data.([]interface{}) {
		v := v.(map[string]interface{})
		id := v["id"].(map[string]interface{})
		snippet := v["snippet"].(map[string]interface{})
		thumbnails := snippet["thumbnails"].(map[string]interface{})
		def := thumbnails["default"].(map[string]interface{})
		value := map[string]interface{}{
			"title":     snippet["title"],
			"desc":      snippet["description"],
			"pub_date":  snippet["publishTime"],
			"thumb_url": def["url"],
			"url":       "https://www.youtube.com/watch?v=" + id["videoId"].(string),
		}
		values = append(values, value)
	}
	json_data, err := json.Marshal(values)
	if err != nil {
		log.Println(err)
	}

	resp, err = http.Post("http://localhost:8080/api/video/", "application/json",
		bytes.NewBuffer(json_data))
	if err != nil {
		log.Println(err)
	}
	log.Println(resp)
}
