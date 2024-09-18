package views

import (
	"context"
	"github.com/Noah-Wilderom/video-streaming-test/models"
)

const GlobalDataContextKey = "global_data"

type GlobalData struct {
	Name        string
	Host        string
	Version     string
	Environment string
	Debug       bool
}

func GetGlobalData(ctx context.Context) GlobalData {
	data, ok := ctx.Value(GlobalDataContextKey).(GlobalData)
	if !ok {
		return GlobalData{}
	}

	return data
}

func AuthenticatedUser(ctx context.Context) models.AuthenticatedUser {
	user, ok := ctx.Value(models.UserContextKey).(models.AuthenticatedUser)
	if !ok {
		return models.AuthenticatedUser{}
	}

	return user
}
