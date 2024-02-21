package service

import (
	"encoding/json"
	"time"

	"github.com/KaySar12/NextZenOS/model"
	"github.com/KaySar12/NextZenOS/pkg/config"
	"github.com/KaySar12/NextZenOS/pkg/utils/httper"
	"github.com/tidwall/gjson"
)

type CasaService interface {
	GetCasaosVersion() model.Version
}

type casaService struct{}

/**
 * @description: get remote version
 * @return {model.Version}
 */

func getLatestVersion() model.Version {
	v := httper.OasisGet(config.ServerInfo.ServerApi + "/v1/sys/version")
	data := gjson.Get(v, "data")
	newVersion := model.Version{}
	err := json.Unmarshal([]byte(data.String()), &newVersion)
	if err != nil {
		panic(err) // Handle error appropriately
	}

	return model.Version{
		Id:        1,
		ChangeLog: newVersion.ChangeLog,
		Version:   "1.1",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
func (o *casaService) GetCasaosVersion() model.Version {
	keyName := "casa_version"
	// var dataStr string
	var version model.Version

	// Check cache and return version if present
	if result, ok := Cache.Get(keyName); ok {
		dataStr, ok := result.(string)
		if ok {
			data := gjson.Parse(dataStr) // Parse as gjson.Result directly
			json.Unmarshal([]byte(data.Get("data").String()), &version)
			return version
		}
	}
	// Directly unmarshal into a new struct instance to avoid pass-by-value issues
	newVersion := getLatestVersion()
	// Cache the modified version
	if len(newVersion.Version) > 0 {
		Cache.Set(keyName, newVersion, time.Minute*20)
	}

	return newVersion
}

func NewCasaService() CasaService {
	return &casaService{}
}
