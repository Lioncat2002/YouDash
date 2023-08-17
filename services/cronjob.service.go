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

// Base Youtube data v3 API url
const BASE_URL = "https://youtube.googleapis.com/youtube/v3/search?part=snippet&order=date&publishedAfter=2023-08-01T00%3A00%3A00Z&q=official&type=video&key="

// For getting the working GCP API KEY from the multiple provided ones
var Count = 0

// Convert the string into a dynamic map
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
	key := keys[Count%len(keys)] //get the current key based on the Count
	log.Println("Using Key: ", key)
	//call youtube api
	resp, err := http.Get(BASE_URL + key)
	log.Println(resp.Status)
	//incase the key has exceeded Quota
	if strings.Contains(resp.Status, "403") {
		Count += 1 //go to next provided key
		log.Println(key, " Quota Exceeded")
	}
	if err != nil {
		log.Println("Error fetching from youtube", err)
		return
	}
	defer resp.Body.Close()
	//convert into byte[]
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	sb := string(body)
	res := jsonToMap(sb)
	data := res["items"] //get the videos datas
	var values []map[string]interface{}
	if data != nil {

		for _, v := range data.([]interface{}) {
			// v.(map[string]interface{})  where (map[string]interface{}) for ensuring type and pleasing go compiler
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
			//turn the data into a dynamic map with schema which can be accepted by the backend
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
		//send the data to the backend as a POST request
		resp, err = http.Post("http://localhost:8080/api/video/", "application/json",
			bytes.NewBuffer(json_data))
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(resp)
	}
}
