package main

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

func staticFile(file string) func(Response, url.Values) ([]byte, string, int) {
	return func(Response, url.Values) ([]byte, string, int) {
		b, _ := ReadFile(file)
		return b, "text/html", 200
	}
}

func apiAll(res Response, _ url.Values) ([]byte, string, int) {
	return []byte(res.ToJson()), "application/json", 200
}

func apiTLS(res Response, _ url.Values) ([]byte, string, int) {
	return []byte(Response{
		TLS: res.TLS,
	}.ToJson()), "application/json", 200
}

func apiClean(res Response, _ url.Values) ([]byte, string, int) {
	akamai := "-"
	hash := "-"
	if res.HTTPVersion == "h2" {
		akamai = res.Http2.AkamaiFingerprint
		hash = GetMD5Hash(res.Http2.AkamaiFingerprint)
	}
	return []byte(SmallResponse{
		JA3:           res.TLS.JA3,
		JA3Hash:       res.TLS.JA3Hash,
		Akamai:        akamai,
		AkamaiHash:    hash,
		PeetPrint:     res.TLS.PeetPrint,
		PeetPrintHash: res.TLS.PeetPrintHash,
	}.ToJson()), "application/json", 200
}

func apiRequestCount(_ Response, _ url.Values) ([]byte, string, int) {
	if !connectedToDB {
		return []byte("{\"error\": \"Not connected to database.\"}"), "application/json", 500
	}
	return []byte(fmt.Sprintf(`{"total_requests": %v}`, GetTotalRequestCount())), "application/json", 200
}

func apiSearchJA3(_ Response, u url.Values) ([]byte, string, int) {
	if !connectedToDB {
		return []byte("{\"error\": \"Not connected to database.\"}"), "application/json", 500
	}
	by := getParam("by", u)
	if by == "" {
		return []byte("{\"error\": \"No 'by' param present\"}"), "application/json", 500
	}
	res := GetByJa3(by)
	j, _ := json.MarshalIndent(res, "", "\t")
	return j, "application/json", 200
}

func apiSearchH2(_ Response, u url.Values) ([]byte, string, int) {
	if !connectedToDB {
		return []byte("{\"error\": \"Not connected to database.\"}"), "application/json", 500
	}
	by := getParam("by", u)
	if by == "" {
		return []byte("{\"error\": \"No 'by' param present\"}"), "application/json", 500
	}
	res := GetByH2(by)
	j, _ := json.MarshalIndent(res, "", "\t")
	return j, "application/json", 200
}

func apiSearchPeetPrint(_ Response, u url.Values) ([]byte, string, int) {
	if !connectedToDB {
		return []byte("{\"error\": \"Not connected to database.\"}"), "application/json", 500
	}
	by := getParam("by", u)
	if by == "" {
		return []byte("{\"error\": \"No 'by' param present\"}"), "application/json", 500
	}
	res := GetByPeetPrint(by)
	j, _ := json.MarshalIndent(res, "", "\t")
	return j, "application/json", 200
}

func apiSearchUserAgent(_ Response, u url.Values) ([]byte, string, int) {
	if !connectedToDB {
		return []byte("{\"error\": \"Not connected to database.\"}"), "application/json", 500
	}
	by := getParam("by", u)
	if by == "" {
		return []byte("{\"error\": \"No 'by' param present\"}"), "application/json", 500
	}
	res := GetByUserAgent(by)
	j, _ := json.MarshalIndent(res, "", "\t")
	return j, "application/json", 200
}

func index(r Response, v url.Values) ([]byte, string, int) {
	res, ct, respCode := staticFile("static/index.html")(r, v)
	data, _ := json.Marshal(r)
	return []byte(strings.ReplaceAll(string(res), "/*DATA*/", string(data))), ct, respCode
}

func getAllPaths() map[string]func(Response, url.Values) ([]byte, string, int) {
	return map[string]func(Response, url.Values) ([]byte, string, int){
		//"/":                     index,
		//"/explore":              staticFile("static/explore.html"),
		"/api/all":   apiAll,
		"/api/tls":   apiTLS,
		"/api/clean": apiClean,
		//"/api/request-count":    apiRequestCount,
		//"/api/search-ja3":       apiSearchJA3,
		//"/api/search-h2":        apiSearchH2,
		//"/api/search-peetprint": apiSearchPeetPrint,
		//"/api/search-useragent": apiSearchUserAgent,
	}
}
