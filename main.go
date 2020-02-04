package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/buaazp/fasthttprouter"
	_ "github.com/lib/pq"
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

type ItemResponse struct {
	Domain  string           `json:"domain"`
	Servers []ServerResponse `json:"servers"`
}

type ItemsResponse struct {
	Items []ItemResponse `json:"items"`
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
			log.Fatal("Error fetching data")
		}

	}

	servers := ServersResponse{}

	for _, endpoint := range ser.Endpoints {
		servers.Servers = append(servers.Servers, ServerResponse{Address: endpoint.IPAddress, SslGrade: endpoint.Grade})
		// println(endpoint.IPAddress)
	}

	InsertHostInfo(ser)

	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	json.NewEncoder(ctx).Encode(servers)

	// json.NewDecoder(bodyBytes).Decode(ctx)

	// println(string(bodyBytes))
}

func InsertEndpointInfo(newEndpoint Endpoint, hostID int) {

	db, err := sql.Open("postgres", "postgresql://root@localhost:26257/domains_db?sslmode=disable")

	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	var returnedID int
	err2 := db.QueryRow(
		`INSERT INTO endpoints_info (address, ssl_grade, domain_id) 
		VALUES ($1, $2, $3) RETURNING endpoint_id;`, newEndpoint.IPAddress, newEndpoint.Grade, hostID).Scan(&returnedID)

	if err2 != nil {
		log.Fatal("Error executing query: ", err2)
	}

}

func InsertHostInfo(newServer Server) {

	db, err := sql.Open("postgres", "postgresql://root@localhost:26257/domains_db?sslmode=disable")

	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}

	var returnedID int
	err2 := db.QueryRow(
		`INSERT INTO domains_info (host_name) 
		VALUES ($1) RETURNING domain_id;`, newServer.Host).Scan(&returnedID)

	if err2 != nil {
		log.Fatal("Error executing query: ", err2)
	}

	for _, endpoint := range newServer.Endpoints {
		InsertEndpointInfo(endpoint, returnedID)
	}

}

func GetHostsInfo(ctx *fasthttp.RequestCtx) {

	db, err := sql.Open("postgres", "postgresql://root@localhost:26257/domains_db?sslmode=disable")

	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	// Select Statement.
	rows, err := db.Query("SELECT domain_id, host_name FROM domains_info;")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	items := ItemsResponse{}

	for rows.Next() {
		var domainID int
		var hostName string

		if err := rows.Scan(&domainID, &hostName); err != nil {
			log.Fatal(err)
		}

		rows2, err := db.Query(`SELECT address, ssl_grade FROM endpoints_info WHERE domain_id = $1;`, domainID)

		if err != nil {
			log.Fatal(err)
		}

		defer rows2.Close()

		item := ItemResponse{Domain: hostName}

		for rows2.Next() {
			var address string
			var ssl_grade string

			if err := rows2.Scan(&address, &ssl_grade); err != nil {
				log.Fatal(err)
			}

			item.Servers = append(item.Servers, ServerResponse{Address: address, SslGrade: ssl_grade})
		}

		items.Items = append(items.Items, item)

	}

	ctx.Response.Header.SetCanonical(strContentType, strApplicationJSON)
	json.NewEncoder(ctx).Encode(items)
}

func InitDB() {
	db, err := sql.Open("postgres", "postgresql://root@localhost:26257/domains_db?sslmode=disable")

	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	if _, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS "domains_info" (
			"domain_id" SERIAL,
			"host_name" STRING(100),
			PRIMARY KEY ("domain_id")
			);`); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec(
		`CREATE TABLE IF NOT EXISTS "endpoints_info" (
				"endpoint_id" SERIAL,
				"address" STRING(100),
				"ssl_grade" STRING(50),
				"domain_id" INT NOT NULL REFERENCES "domains_info" (domain_id) ON DELETE CASCADE,
				PRIMARY KEY ("endpoint_id")
		);`); err != nil {
		log.Fatal(err)
	}
}

func main() {

	InitDB()

	router := fasthttprouter.New()
	router.GET("/", Index)
	router.GET("/servers/:host", ServerInfo)
	router.GET("/history", GetHostsInfo)

	log.Fatal(fasthttp.ListenAndServe(":8081", router.Handler))
}
