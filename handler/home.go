package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Noah-Wilderom/video-streaming-test/proto/video"
	"github.com/Noah-Wilderom/video-streaming-test/resources/views/home"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
)

func HandleHomeIndex(c echo.Context) error {
	user := getAuthenticatedUser(c)
	apiUrl := "http://localhost:8080/video"
	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", user.Token))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(string(respBody))

	var respStruct struct {
		Videos []*video.Video `json:"videos"`
	}

	err = json.Unmarshal(respBody, &respStruct)

	return home.Index(respStruct.Videos).Render(c.Request().Context(), c.Response().Writer)
}
