package dropbox

import (
	"github.com/IceWhaleTech/CasaOS/internal/driver"
)

const ICONURL = "./img/driver/Dropbox.svg"

type Addition struct {
	driver.RootID
	RefreshToken   string `json:"refresh_token" required:"true" omit:"true"`
	AppKey         string `json:"app_key" type:"string" default:"y5dn1dhp6yh4vwe" omit:"true"`
	AppSecret      string `json:"app_secret" type:"string" default:"mhlkkb1axzxpy4r" omit:"true"`
	OrderDirection string `json:"order_direction" type:"select" options:"asc,desc" omit:"true"`
	AuthUrl        string `json:"auth_url" type:"string" default:""`
	Icon           string `json:"icon" type:"string" default:"./img/driver/Dropbox.svg"`
	Code           string `json:"code" type:"string" help:"code from auth_url" omit:"true"`
}

var config = driver.Config{
	Name:        "Dropbox",
	OnlyProxy:   true,
	DefaultRoot: "root",
}
