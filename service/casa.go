package service

import (
	json "encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/tidwall/gjson"
)

type CasaService interface {
	GetCasaosVersion() (model.Version, error)
}

type casaService struct{}

func getLatestVersion() (model.Version, error) {
	resp, err := http.Get("https://api.nextzenos.com/v1/sys/version")
	if err != nil {
		return model.Version{}, fmt.Errorf("failed to fetch latest version: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return model.Version{}, fmt.Errorf("failed to read response body: %w", err)
	}

	data := gjson.Get(string(body), "data")
	newVersion := model.Version{}

	err = json.Unmarshal([]byte(data.String()), &newVersion)
	if err != nil {
		return model.Version{}, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return newVersion, nil
}

func (o *casaService) GetCasaosVersion() (model.Version, error) {
	keyName := "casa_version"
	var version model.Version

	// Check cache and return version if present
	if result, ok := Cache.Get(keyName); ok {
		dataStr, ok := result.(string)
		if ok {
			data := gjson.Parse(dataStr) // Parse as gjson.Result directly
			json.Unmarshal([]byte(data.Get("data").String()), &version)
			return version, nil
		}
	}
	// Fetch latest version if not in cache
	newVersion, err := getLatestVersion()
	if err != nil {
		return model.Version{}, fmt.Errorf("failed to fetch latest version: %w", err)
	}

	// Cache the fetched version
	if len(newVersion.Version) > 0 {
		Cache.Set(keyName, newVersion, time.Minute*20)
	}

	return newVersion, nil
}

func NewCasaService() CasaService {
	return &casaService{}
}
