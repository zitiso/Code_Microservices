package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

//error handling for all errors! Not very graceful but sufficient for the lab
//better to set up error responses to tell client what happened and respond with appropriate status code
func check(e error) {
	if e != nil {
		panic(e)
	}
}

//handler for the info request - responds with text
func infoRequest_Handler(c *gin.Context) {

	//send text response with success status code
	c.String(200, "Welcome to the Weather App microservice.")
}

//handler for the weather request - responds with json
func weatherRequest_Handler(c *gin.Context) {

	// backend service properties
	service_url := os.Getenv("WEATHER_API_URL")
	apikey := "&appid=" + os.Getenv("WEATHER_API_KEY")

	//request parameters to use when calling the weather service
	//others could be added - optional ones could have a default set here
	unit := "&units=metric"
	city := "?q=" + c.Query("city")

	//the request url needs to be created dynamically
	ask_url := service_url + city + apikey + unit
	log.Println("Ask requested: ", ask_url)

	//invoke backend service and check for errors
	resp, err := http.Get(ask_url)
	check(err)

	//process the backend service response and check for errors
	body, err := ioutil.ReadAll(resp.Body)
	check(err)

	//prepare the client response and check for errors
	w := map[string]interface{}{}
	err1 := json.Unmarshal(body, &w)
	check(err1)
	// log.Println("Unmarshaled json: ", w)

	//send JSON response with success status code
	c.JSON(200, w)

	//Update report using Docker Volume
    //If not in Docker, create file on file system and update code to use it
	writeVolume("Weather for " + c.Query("city") + " was just sent")
}

func writeVolume(m string) {
	//check that volume mounted and file exists
	_, err := os.Stat("/data/weather/report.txt")
	if err != nil {
		fmt.Println("Volume not monted ")
		os.Exit(0)
	}

	//read in file mounted from volume
	existing, err := os.Open("/data/weather/report.txt")
	check(err)
	defer existing.Close()

	scanner := bufio.NewScanner(existing)

	fmt.Println("before for loop")
	for scanner.Scan() {
		fmt.Println("Text before change :", scanner.Text())
		fmt.Println("in reading for loop")
	}

	//write out content to share and persist using the volume
	request_to_share, err2 := os.OpenFile("/data/weather/report.txt", os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	check(err2)
	defer request_to_share.Close()

	content := []byte(m + "\n")
	len, err := request_to_share.Write(content)
	check(err)

	fmt.Printf("\nLength: %d bytes \n", len)
	fmt.Printf("\nFile name: %s \n", request_to_share.Name())

	//verify writing out worked
	requests, err := os.Open("/data/weather/report.txt")
	check(err)
	defer requests.Close()

	scanner2 := bufio.NewScanner(requests)

	fmt.Println("before for loop")
	size := 0
	for scanner2.Scan() {
		fmt.Println("Text after change :", scanner2.Text())
		fmt.Println("in verification for loop")
		size += 1
	}
	fmt.Printf("Number of entries: %d\n", size)

}

func main() {
	//set up routes
	r := gin.Default()

	r.GET("/info", infoRequest_Handler)

	r.GET("/weather", weatherRequest_Handler)

	//start server listening on port specified in the env_file
	r.Run(os.Getenv("PORT"))
}
