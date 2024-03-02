package onedrive

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/IceWhaleTech/CasaOS-Common/utils/logger"
	"github.com/IceWhaleTech/CasaOS/drivers/base"
	"go.uber.org/zap"
)

var (
	client_id     = "b23a17d9-4d9d-4fad-a0a9-7c905a6e3fe7"
	client_secret = "bx~8Q~_VmPBAu24JcFSDGrIozTGQHFRbrRIjmaWy"
)
var onedriveHostMap = map[string]Host{
	"global": {
		Oauth: "https://login.microsoftonline.com",
		Api:   "https://graph.microsoft.com",
	},
	"cn": {
		Oauth: "https://login.chinacloudapi.cn",
		Api:   "https://microsoftgraph.chinacloudapi.cn",
	},
	"us": {
		Oauth: "https://login.microsoftonline.us",
		Api:   "https://graph.microsoft.us",
	},
	"de": {
		Oauth: "https://login.microsoftonline.de",
		Api:   "https://graph.microsoft.de",
	},
}

func EncodePath(path string, all ...bool) string {
	seg := strings.Split(path, "/")
	toReplace := []struct {
		Src string
		Dst string
	}{
		{Src: "%", Dst: "%25"},
		{"%", "%25"},
		{"?", "%3F"},
		{"#", "%23"},
	}
	for i := range seg {
		if len(all) > 0 && all[0] {
			seg[i] = url.PathEscape(seg[i])
		} else {
			for j := range toReplace {
				seg[i] = strings.ReplaceAll(seg[i], toReplace[j].Src, toReplace[j].Dst)
			}
		}
	}
	return strings.Join(seg, "/")
}
func (d *Onedrive) GetMetaUrl(auth bool, path string) string {
	host := onedriveHostMap[d.Region]
	path = EncodePath(path, true)
	if auth {
		return host.Oauth
	}
	if d.IsSharepoint {
		if path == "/" || path == "\\" {
			return fmt.Sprintf("%s/v1.0/sites/%s/drive/root", host.Api, d.SiteId)
		} else {
			return fmt.Sprintf("%s/v1.0/sites/%s/drive/root:%s:", host.Api, d.SiteId, path)
		}
	} else {
		if path == "/" || path == "\\" {
			return fmt.Sprintf("%s/v1.0/me/drive/root", host.Api)
		} else {
			return fmt.Sprintf("%s/v1.0/me/drive/root:%s:", host.Api, path)
		}
	}
}

func (d *Onedrive) refreshToken() error {
	var err error
	for i := 0; i < 3; i++ {
		err = d._refreshToken()
		if err == nil {
			break
		}
	}
	return err
}

func (d *Onedrive) getRefreshToken() error {
	url := d.GetMetaUrl(true, "") + "/common/oauth2/v2.0/token"
	var resp base.TokenResp
	var e TokenErr

	res, err := base.RestyClient.R().SetResult(&resp).SetError(&e).SetFormData(map[string]string{
		"grant_type":    "authorization_code",
		"client_id":     d.ClientID,
		"client_secret": d.ClientSecret,
		"code":          d.Code,
		"redirect_uri":  d.RedirectUri,
	}).Post(url)
	if err != nil {
		return err
	}
	logger.Info("get refresh token", zap.String("res", res.String()))
	if e.Error != "" {
		return fmt.Errorf("%s", e.ErrorDescription)
	}
	if resp.RefreshToken == "" {
		return errors.New("refresh token is empty")
	}
	d.RefreshToken, d.AccessToken = resp.RefreshToken, resp.AccessToken
	return nil
}

func (d *Onedrive) _refreshToken() error {
	url := d.GetMetaUrl(true, "") + "/common/oauth2/v2.0/token"
	var resp base.TokenResp
	var e TokenErr

	res, err := base.RestyClient.R().SetResult(&resp).SetError(&e).SetFormData(map[string]string{
		"grant_type":    "refresh_token",
		"client_id":     d.ClientID,
		"client_secret": d.ClientSecret,
		"redirect_uri":  d.RedirectUri,
		"refresh_token": d.RefreshToken,
	}).Post(url)
	if err != nil {
		return err
	}
	logger.Info("get refresh token", zap.String("res", res.String()))
	if e.Error != "" {
		return fmt.Errorf("%s", e.ErrorDescription)
	}
	if resp.RefreshToken == "" {
		return errors.New("refresh token is empty")
	}
	d.RefreshToken, d.AccessToken = resp.RefreshToken, resp.AccessToken
	return nil
}

func (d *Onedrive) Request(url string, method string, callback base.ReqCallback, resp interface{}) ([]byte, error) {
	req := base.RestyClient.R()
	req.SetHeader("Authorization", "Bearer "+d.AccessToken)
	if callback != nil {
		callback(req)
	}
	if resp != nil {
		req.SetResult(resp)
	}
	var e RespErr
	req.SetError(&e)
	res, err := req.Execute(method, url)
	if err != nil {
		return nil, err
	}
	if e.Error.Code != "" {
		if e.Error.Code == "InvalidAuthenticationToken" {
			err = d.refreshToken()
			if err != nil {
				return nil, err
			}
			return d.Request(url, method, callback, resp)
		}
		return nil, errors.New(e.Error.Message)
	}
	return res.Body(), nil
}

func GetConfig() Onedrive {
	config := Onedrive{}
	config.ClientID = client_id
	config.ClientSecret = client_secret
	config.RootFolderID = "/"
	config.AuthUrl = "https://login.microsoftonline.com/common/oauth2/v2.0/authorize?client_id=" + client_id + "&response_type=code&redirect_uri=https%3A%2F%2Fnextzen-api.onrender.com/cloud&scope=offline_access+files.readwrite.all&state=${HOST}%2Fv1%2Frecover%2FOnedrive"
	config.Icon = "./img/driver/OneDrive.svg"
	config.Region = "global"
	config.RedirectUri = "https://nextzen-api.onrender.com/cloud"

	return config
}
