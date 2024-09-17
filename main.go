package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/Noah-Wilderom/video-streaming-test/resources/views"
	"github.com/labstack/echo/v4"
)

const (
	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjcmVhdGVkX2F0IjoiMjAyNC0wOS0xN1QxNjo0OTowMC43MDY1MDQ0NTdaIiwiZGF0YSI6eyJ1c2VyIjoiR1Q1UlZUWmhaZ0k5UGxDNXdUaEpSVzV6SmpkTGZITk5VZzhIS1YxMGRqSkdjaEVLSmg0b2UzMUhSR0l5V0FjQU5RQXFYVkJIVEhGUkN3ZFRXRUFLVXc0MVU2dTFaVEo0RFR3MUp6VmRMQUliY0pEUlNuQnpYS25BWTBVNWNTNWJLaEpWbU1WM1JIcU85N3JxWXgxRVFudGVQZzVqZm1GRFBYRkJHdzRGV1ZOaWVXWXlCMGhTUjJnb00zY2pDU2x3WDBWdVhnUWZFaVIyRmc4b1hCQUlYU3AyUWhBdUd5QWZBeWdvQkExWUhETWxBamxRUFdVdkh6VkpTUklMVlVzSGUyVWlHVk1VWmhVL2VqTWdLMWN3QWlRakFIMHBJUlVwRGpjbkVYQTVPMmd3QW1FamZnUUNLaFlOY3g1WGJpUXhNaFVMYWdNOEIzRlhTWFUyT1RxeFRtVExVbFEzWlROeFdVTlFSa0U0UnAwV1d1VlRTbkJ6VEZaQiJ9LCJleHAiOjE3MjY1OTUzNDB9.MxxcLpjU_Tj36r1zIoK2pO219L17IpVtBPEDd6zQKCc"
)

func main() {
	e := echo.New()

	e.StaticFS("/public/*", os.DirFS("public"))

	e.GET("/", func(c echo.Context) error {
		return views.Index().Render(c.Request().Context(), c.Response().Writer)
	})

	e.GET("/video/:id", func(c echo.Context) error {
		videoId := c.Param("id")
		streamApiUrl := "http://localhost:8080/stream/new/" + videoId
		streamReq, err := http.NewRequest(http.MethodPost, streamApiUrl, nil)
		streamReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
		streamClient := &http.Client{}
		streamResp, err := streamClient.Do(streamReq)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer streamResp.Body.Close()

		if streamResp.StatusCode != http.StatusOK {
			return c.String(http.StatusNotFound, "video not found")
		}

		streamRespBody, err := io.ReadAll(streamResp.Body)
		if err != nil {
			fmt.Println(err)
			return err
		}

		var streamRespStruct struct {
			Content string `json:"content"`
			Type    string `json:"type"`
		}

		err = json.Unmarshal(streamRespBody, &streamRespStruct)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return views.Video(videoId, true, streamRespStruct.Content, token).Render(c.Request().Context(), c.Response().Writer)

	})

	e.POST("/upload", func(c echo.Context) error {
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
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
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

		return views.UploadInProgress("123").Render(c.Request().Context(), c.Response().Writer)
	})

	e.GET("/upload/progress/:id", func(c echo.Context) error {
		id := c.QueryParam("id")

		return views.UploadInProgress(id).Render(c.Request().Context(), c.Response().Writer)
	})

	e.Logger.Fatal(e.Start(":3000"))
}
