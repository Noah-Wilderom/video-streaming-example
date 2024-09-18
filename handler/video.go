package handler

import (
	"encoding/json"
	"fmt"
	"github.com/Noah-Wilderom/video-streaming-test/resources/views/home"
	"github.com/Noah-Wilderom/video-streaming-test/resources/views/video"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"strconv"
)

func HandleVideoShow(c echo.Context) error {
	user := getAuthenticatedUser(c)

	videoId := c.Param("id")
	streamApiUrl := "http://localhost:8080/stream/new/" + videoId
	streamReq, err := http.NewRequest(http.MethodPost, streamApiUrl, nil)
	streamReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", user.Token))
	streamClient := &http.Client{}
	streamResp, err := streamClient.Do(streamReq)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer streamResp.Body.Close()

	if streamResp.StatusCode != http.StatusOK {
		fmt.Println(streamResp.StatusCode)
		if streamResp.StatusCode == http.StatusNotFound {
			return c.String(http.StatusNotFound, "video not found")
		}
	}

	streamRespBody, err := io.ReadAll(streamResp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	var streamRespStruct struct {
		Content string `json:"content"`
		Type    string `json:"type"`
	}

	fmt.Println(string(streamRespBody))
	err = json.Unmarshal(streamRespBody, &streamRespStruct)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return Render(c, video.Video(videoId, true, streamRespStruct.Content, user.Token))
}

func HandleUploadVideo(c echo.Context) error {
	user := getAuthenticatedUser(c)
	file, header, err := c.Request().FormFile("file")
	if err != nil {
		return err
	}

	mimeType := header.Header.Get("Content-Type")
	fileSize := header.Size
	go func() {
		defer file.Close()
		// Prepare the request to forward the file to the external API
		apiURL := "http://localhost:8080/video/upload"
		req, err := http.NewRequest(http.MethodPost, apiURL, file)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Filesize:", strconv.FormatInt(fileSize, 10))

		// Set the Content-Length header with the file size
		req.Header.Set("Content-Length", strconv.FormatInt(fileSize, 10))
		req.Header.Set("Content-Type", "application/octet-stream") // Assuming the external API expects raw data
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", user.Token))
		req.Header.Set("Content-Mimetype", mimeType)
		req.ContentLength = fileSize

		// Perform the request to the external API
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()

		// Read the response from the external API
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
		}

		// Return success or error based on the API response
		if resp.StatusCode != http.StatusOK {
			fmt.Println(resp.StatusCode)
			fmt.Println(string(respBody))
		}
	}()

	return Render(c, home.UploadInProgress("123"))
}
