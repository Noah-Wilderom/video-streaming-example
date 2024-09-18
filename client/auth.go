package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

func (c *Client) Login(ctx context.Context, email, password string) (*LoginResponse, int, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	_ = writer.WriteField("email", email)
	_ = writer.WriteField("password", password)

	err := writer.Close()
	if err != nil {
		fmt.Println("Error closing writer:", err)
		return nil, 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/auth/login", c.BaseUrl), &body)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Message", string(respBody))
		return nil, resp.StatusCode, fmt.Errorf("login error: %s", resp.Status)
	}

	var loginResponse *LoginResponse
	err = json.Unmarshal(respBody, &loginResponse)
	if err != nil {
		fmt.Println("resp body", string(respBody))
		fmt.Println(err.Error())
		return nil, 0, err
	}

	return loginResponse, resp.StatusCode, nil
}

func (c *Client) Signup(ctx context.Context, name, email, password string) (*LoginResponse, int, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	_ = writer.WriteField("name", name)
	_ = writer.WriteField("email", email)
	_ = writer.WriteField("password", password)

	err := writer.Close()
	if err != nil {
		fmt.Println("Error closing writer:", err)
		return nil, 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/auth/signup", c.BaseUrl), &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(string(respBody))
		return nil, resp.StatusCode, fmt.Errorf("signup error: %s", resp.Status)
	}

	var loginResponse *LoginResponse
	err = json.Unmarshal(respBody, &loginResponse)
	if err != nil {
		return nil, 0, err
	}

	return loginResponse, resp.StatusCode, nil
}

func (c *Client) Check(ctx context.Context, token string) (*LoginResponse, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	err := writer.Close()
	if err != nil {
		fmt.Println("Error closing writer:", err)
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/auth/check", c.BaseUrl), &body)

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("check", string(respBody))
		return nil, fmt.Errorf("check error: %s", resp.Status)
	}

	var loginResponse *LoginResponse
	err = json.Unmarshal(respBody, &loginResponse)
	if err != nil {
		return nil, err
	}

	return loginResponse, nil
}
