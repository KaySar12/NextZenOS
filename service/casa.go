package service

import (
	json "encoding/json"
	"time"

	"github.com/IceWhaleTech/CasaOS/model"
	"github.com/IceWhaleTech/CasaOS/pkg/utils/httper"
	"github.com/tidwall/gjson"
)

type CasaService interface {
	GetCasaosVersion() model.Version
}

type casaService struct{}

func getLatestVersion() model.Version {
	v := httper.OasisGet("https://nextzen-api.onrender.com" + "/v1/sys/version")
	data := gjson.Get(v, "data")
	newVersion := model.Version{}
	err := json.Unmarshal([]byte(data.String()), &newVersion)
	if err != nil {
		panic(err) // Handle error appropriately
	}

	return model.Version{
		Id:        1,
		ChangeLog: newVersion.ChangeLog,
		Version:   newVersion.Version,
		CreatedAt: newVersion.CreatedAt,
		UpdatedAt: newVersion.UpdatedAt,
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
