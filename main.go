package main

import (
	"encoding/json"
	"errors"
	"fmt"
	goku_plugin "github.com/eolinker/goku-plugin"
	"reflect"
	"strings"

)

var builder = new(goku_APIKeyAuthFactory)

func Builder() goku_plugin.PluginFactory {
	return builder
}

type APIKeyNode struct {
	APIKey         string `json:"Apikey"`
	HideCredential bool   `json:"hideCredential"`
	Remark         string `json:"remark"`
}

type APIKeyConf []APIKeyNode

type goku_APIKeyAuth struct {
	conf APIKeyConf
}
type goku_APIKeyAuthFactory struct {
}

func (f *goku_APIKeyAuthFactory) Create(config string, clusterName string, updateTag string, strategyId string, apiId int) (*goku_plugin.PluginObj, error) {

	if config != "" {
		var conf APIKeyConf

		if err := json.Unmarshal([]byte(config), &conf); err != nil {
			//解析配置信息失败
			return nil, fmt.Errorf("[apikey_auth] Parsing plugin config :%s", err.Error())
		}

		return &goku_plugin.PluginObj{
			Access: &goku_APIKeyAuth{
				conf: conf,
			},
		}, nil
	}
	return nil, errors.New("need config")
}

func (p *goku_APIKeyAuth) Access(ctx goku_plugin.ContextAccess) (isContinue bool, e error) {



	for _, node := range p.conf {
		_, flag, err := retrieveAPIKeyCredential(ctx, node)
		if flag {
			//认证成功
			return true, nil
		}
		if err != nil {
			ctx.SetStatus(403, "403")

			ctx.SetBody([]byte(err.Error()))
			return false, err
		}
	}
	//认证信息不匹配
	ctx.SetStatus(403, "403")

	ctx.SetBody([]byte("[apikey_auth] Invalid Apikey"))
	return false, errors.New("[apikey_auth] Invalid Apikey")

}

//获取验证信息
func retrieveAPIKeyCredential(ctx goku_plugin.ContextAccess, node APIKeyNode) (string, bool, error) {
	v := ""
	if values, ok := ctx.Request().Headers()["Authorization"]; ok {
		v = values[0]
		if v != "" && v == node.APIKey {
			if node.HideCredential {
				ctx.Proxy().DelHeader("Authorization")
			}
			return v, true, nil
		}
		return v, false, nil
	}
	if values, ok := ctx.Request().Headers()["Apikey"]; ok {
		v := values[0]
		if v != "" && v == node.APIKey {
			if node.HideCredential {
				ctx.Proxy().DelHeader("Apikey")
			}
			return v, true, nil
		}
		return v, false, nil
	}

	if values, ok := ctx.Request().URL().Query()["Apikey"]; ok {
		v = values[0]
		if v != "" && v == node.APIKey {
			if node.HideCredential {
				ctx.Proxy().Querys().Del("Apikey")
			}
			return v, true, nil
		}
		return v, false, nil
	}

	contentType := ctx.Request().Headers().Get("Content-Type")
	if strings.Contains(contentType, "application/x-www-form-urlencoded") || strings.Contains(contentType, "application/www-form-urlencoded") {

		formParams, err := ctx.Proxy().BodyForm()
		if err != nil {
			return "", false, err
		}
		if _, ok := formParams["Apikey"]; ok {
			v = formParams["Apikey"][0]
			if v != "" && v == node.APIKey {
				if node.HideCredential {
					delete(formParams, "Apikey")
					ctx.Proxy().SetForm(formParams)
				}
				return v, true, nil
			}
			return v, false, nil
		}
	} else if strings.Contains(contentType, "application/json") {
		var body map[string]interface{}


		rawbody, err := ctx.Proxy().RawBody()
		if err != nil {
			return "", false, err
		}
		if err := json.Unmarshal(rawbody, &body); err != nil {
			return "", false, err
		}
		if _, ok := body["Apikey"]; !ok {
			return "", false, errors.New("[apikey_auth] cannot find the Apikey in body")
		}
		if TOfData(body["Apikey"]) == reflect.String {
			v = body["Apikey"].(string)
			if v != "" && v == node.APIKey {
				if node.HideCredential {
					delete(body, "Apikey")
					editedBody, err := json.Marshal(body)
					if err != nil {
						return "", false, err
					}
					ctx.Proxy().SetRaw(contentType, editedBody)
				}
				return v, true, nil
			}
			return v, false, nil
		} else {
			return "", false, errors.New("[apikey_auth] Invalid data type for Apikey")
		}
	} else if strings.Contains(contentType, "multipart/form-data") {
		bodyform, err := ctx.Proxy().BodyForm()
		if err != nil {
			return "", false, err
		}
		if values, ok := bodyform["Apikey"]; ok {
			v = values[0]
			if v != "" && v == node.APIKey {
				if node.HideCredential {
					delete(bodyform, "Apikey")
					ctx.Proxy().SetForm(bodyform)
				}

				return v, true, nil
			}
			return v, false, nil
		}
	} else {
		return "", false, errors.New("[apikey_auth] Unsupported Content-Type")
	}
	if v == "" {
		return "", false, errors.New("[apikey_auth] cannot find the Apikey in query/body/header")
	}
	return v, true, nil
}
