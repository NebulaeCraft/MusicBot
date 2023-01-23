package Bili

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func QueryVideoTitle(serial string) (string, error) {
	url := "https://www.bilibili.com/video/" + serial
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Failed to create request")
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send request")
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response")
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to get music url")
		return "", err
	}
	title := string(body)
	title = title[strings.Index(title, "<h1 title=\"")+11:]
	title = title[:strings.Index(title, "\" class=\"video-title tit\">")]
	title = strings.Replace(title, "/", "", -1)
	return title, nil
}
