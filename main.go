package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

var (
	strContentType     = []byte("Content-Type")
	strApplicationJSON = []byte("application/json")
)

type Endpoint struct {
	IPAddress string `json:"ipAddress"`
	Grade     string `json:"grade"`
}

type Server struct {
	Host      string     `json:"host"`
	Status    string     `json:"status"`
	Endpoints []Endpoint `json:"endpoints"`
}

type ServerResponse struct {
	Address  string `json:"address"`
	SslGrade string `json:"ssl_grade"`
}

type ServersResponse struct {
	Servers []ServerResponse `json:"servers"`
}

// Index .
func Index(ctx *fasthttp.RequestCtx) {
	fmt.Fprint(ctx, "Welcome!\n")
}

// ServerInfo .
func ServerInfo(ctx *fasthttp.RequestCtx) {

	apiURL := fmt.Sprintf("https://api.ssllabs.com/api/v3/analyze?host=%s", ctx.UserValue("host"))

	var ser Server

	for ser.Status != "READY" {
		println("Trying API call...")

		req := fasthttp.AcquireRequest()
		req.SetRequestURI(apiURL)

		resp := fasthttp.AcquireResponse()
		client := &fasthttp.Client{}
		client.Do(req, resp)

		bodyBytes := resp.Body()

		json.Unmarshal(bodyBytes, &ser)

		// Wait 5 or 10 seconds until next API call (I'm using SSLLabs suggested times)
		if ser.Status == "DNS" {
			time.Sleep(5000 * time.Millisecond)
		} else if ser.Status == "IN_PROGRESS" {
			time.Sleep(10000 * time.Millisecond)
		} else if ser.Status == "ERROR" {
			// Handle error
		}

	}

	servers := ServersResponse{}

	for _, endpoint := range ser.Endpoints {
		servers.Servers = append(servers.Servers, ServerResponse{Address: endpoint.IPAddress, SslGrade: endpoint.Grade})
		// println(endpoint.IPAddress)
	}

	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	json.NewEncoder(ctx).Encode(servers)

	// json.NewDecoder(bodyBytes).Decode(ctx)

	// println(string(bodyBytes))
}

func main() {
	router := fasthttprouter.New()
	router.GET("/", Index)
	router.GET("/servers/:host", ServerInfo)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}
