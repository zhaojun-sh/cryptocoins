package rpcutils

import (
	"io/ioutil"
	"net/http"
	neturl "net/url"
	//"strings"
)

func HttpGet(host string, api string, params map[string][]string) ([]byte, error) {
	requrl := host + "/" + api
	values := neturl.Values(params)
	if params != nil {
		requrl = requrl+"?"+values.Encode()
	}
	resp, err := http.Get(requrl)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

