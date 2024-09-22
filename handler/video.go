package handler

import (
	"context"
	"fmt"
	"github.com/Noah-Wilderom/video-streaming-test/client"
	"github.com/Noah-Wilderom/video-streaming-test/resources/views/home"
	"github.com/Noah-Wilderom/video-streaming-test/resources/views/video"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"strconv"
	"time"
)

func HandleVideoShow(c echo.Context) error {
	user := getAuthenticatedUser(c)

	videoId := c.Param("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	apiClient := client.NewClient()
	streamResp, err := apiClient.WatchStream(ctx, videoId, user.Token)
	if err != nil {
		return err
	}

	previewResp, err := apiClient.PreviewStream(ctx, videoId, user.Token)
	if err != nil {
		return err
	}

	return Render(c, video.Video(videoId, true, streamResp.Content, previewResp.Content, user.Token))
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
