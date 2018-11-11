package main

import (
	"log"
	"net/http"
	"os"
)

func getHomework() {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://up.app.astea.net/tasks/2/export", nil)
	if err != nil {
		log.Fatalf("Error %s", err)
	}
	req.SetBasicAuth(os.Getenv("UP_ASTEA_USER"), os.Getenv("UP_ASTEA_PASSWORD"))
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Error in response %s \n", err)
	}
	log.Printf("Status %s\n", resp.Status)
	defer resp.Body.Close()
	log.Printf("Length of response is %d\n", resp.ContentLength)
}
