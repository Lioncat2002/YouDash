package services

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const BASE_URL = "https://youtube.googleapis.com/youtube/v3/search?part=snippet&order=date&publishedAfter=2023-08-01T00%3A00%3A00Z&q=official&type=video&key="

var Count = 0

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
	keys := strings.Split(os.Getenv("GCP_APIKEY"), " ")
	key := keys[Count%len(keys)]
	log.Println("Using Key: ", key)
	resp, err := http.Get(BASE_URL + key)
	log.Println(resp.Status)
	if strings.Contains(resp.Status, "403") {
		Count += 1
		log.Println(key, " Quota Exceeded")
	}
	if err != nil {
		log.Println("Error fetching from youtube", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	sb := string(body)
	res := jsonToMap(sb)
	data := res["items"]
	var values []map[string]interface{}
	if data != nil {

		for _, v := range data.([]interface{}) {
			v := v.(map[string]interface{})
			id := v["id"].(map[string]interface{})
			snippet := v["snippet"].(map[string]interface{})
			thumbnails := snippet["thumbnails"].(map[string]interface{})
			def := thumbnails["default"].(map[string]interface{})
			resp, err := http.Get("http://localhost:8080/api/video/" + id["videoId"].(string))
			if err != nil {
				log.Println("Error fetching from backend", err)
				return
			}
			if !strings.Contains(resp.Status, "400") {
				log.Println("Skipping ", id["videoId"])
				continue
			}
			value := map[string]interface{}{
				"title":     snippet["title"],
				"desc":      snippet["description"],
				"pub_date":  snippet["publishTime"],
				"thumb_url": def["url"],
				"video_id":  id["videoId"],
			}
			values = append(values, value)
		}
		json_data, err := json.Marshal(values)
		if err != nil {
			log.Println(err)
			return
		}

		resp, err = http.Post("http://localhost:8080/api/video/", "application/json",
			bytes.NewBuffer(json_data))
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(resp)
	}
}
