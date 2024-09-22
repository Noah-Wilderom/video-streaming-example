package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type WatchStreamResponse struct {
	Content string `json:"content"`
	Type    string `json:"type"`
}

func (c *Client) WatchStream(ctx context.Context, videoId string, token string) (*WatchStreamResponse, error) {
	streamApiUrl := "http://localhost:8080/stream/new/" + videoId
	streamReq, err := http.NewRequest(http.MethodPost, streamApiUrl, nil)
	streamReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	streamClient := &http.Client{}
	streamResp, err := streamClient.Do(streamReq)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer streamResp.Body.Close()

	streamRespBody, err := io.ReadAll(streamResp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if streamResp.StatusCode != http.StatusOK {
		fmt.Println(streamResp.StatusCode)
		fmt.Println(streamResp.Body)
		return nil, fmt.Errorf("Video API error [%d]: %s\n", streamResp.StatusCode, streamRespBody)
	}

	var streamResponse WatchStreamResponse
	err = json.Unmarshal(streamRespBody, &streamResponse)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	decrypter := &AES256Encryption{
		Key: []byte(os.Getenv("STREAM_KEY")),
	}
	decryptedContent, err := decrypter.Decrypt(streamResponse.Content)
	if err != nil {
		return nil, err
	}

	streamResponse.Content = decryptedContent
	return &streamResponse, nil
}

func (c *Client) PreviewStream(ctx context.Context, videoId string, token string) (*WatchStreamResponse, error) {
	streamApiUrl := "http://localhost:8080/stream/new-preview/" + videoId
	streamReq, err := http.NewRequest(http.MethodPost, streamApiUrl, nil)
	streamReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	streamClient := &http.Client{}
	streamResp, err := streamClient.Do(streamReq)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer streamResp.Body.Close()

	streamRespBody, err := io.ReadAll(streamResp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	if streamResp.StatusCode != http.StatusOK {
		fmt.Println(streamResp.StatusCode)
		fmt.Println(streamResp.Body)
		return nil, fmt.Errorf("Video API error [%d]: %s\n", streamResp.StatusCode, streamRespBody)
	}

	var streamResponse WatchStreamResponse
	err = json.Unmarshal(streamRespBody, &streamResponse)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	decrypter := &AES256Encryption{
		Key: []byte(os.Getenv("STREAM_KEY")),
	}
	decryptedContent, err := decrypter.Decrypt(streamResponse.Content)
	if err != nil {
		return nil, err
	}

	streamResponse.Content = decryptedContent
	return &streamResponse, nil
}
