// // package main

// // import (
// // 	"encoding/json"
// // 	"io"
// // 	"io/ioutil"
// // 	"log"
// // 	"net/http"
// // 	"os"
// // 	"strconv"

// // 	driver "github.com/johnnadratowski/golang-neo4j-bolt-driver"
// // )

// declare interface MovieResult {}

// declare interface Movie {
//   released: number;
//   title?: string;
//   tagline?: string;
//   cast?: Person[];
// }

// declare interface Person {
//   job: string;
//   role: string[];
//   name: string;
// }

// declare interface D3Response {
//   nodes: Node[];
//   links: Link[];
// }

// declare interface Node {
//   title: string;
//   label: string;
// }

// declare interface Link {
//   source: number;
//   target: number;
// }

// var neo4jURL = "bolt://localhost:7687";

// // function interfaceSliceToString(s []) []string {
// // o := make([]string, len(s))
// // for idx, item := range s {
// // 	o[idx] = item.(string)
// // }
// // return o
// // }

// function defaultHandler(request, response) {
//   // w.Header().Set("Content-Type", "text/html")
//   // body, _ := ioutil.ReadFile("public/index.html")
//   // w.Write(body)
// }

// function searchHandler(response, request) {
// 	// response.Header().Set("Content-Type", "application/json")

// 	// query := request.URL.Query()["q"][0]
// 	// cypher := `
// 	// MATCH
// 	// 	(movie:Movie)
// 	// WHERE
// 	// 	movie.title =~ {query}
// 	// RETURN
// 	// 	movie.title as title, movie.tagline as tagline, movie.released as released`

// 	// db, err := driver.NewDriver().OpenNeo(neo4jURL)
// 	// if err != nil {
// 	// 	log.Println("error connecting to neo4j:", err)
// 	// 	response.WriteHeader(500)
// 	// 	response.Write([]byte("An error occurred connecting to the DB"))
// 	// 	return
// 	// }
// 	// defer db.Close()

// 	// param := "(?i).*" + query + ".*"
// 	// data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{"query": param})
// 	// if err != nil {
// 	// 	log.Println("error querying search:", err)
// 	// 	response.WriteHeader(500)
// 	// 	response.Write([]byte("An error occurred querying the DB"))
// 	// 	return
// 	// } else if len(data) == 0 {
// 	// 	response.WriteHeader(404)
// 	// 	return
// 	// }

// 	// results := make([]MovieResult, len(data))
// 	// for idx, row := range data {
// 	// 	results[idx] = MovieResult{
// 	// 		Movie{
// 	// 			Title:    row[0].(string),
// 	// 			Tagline:  row[1].(string),
// 	// 			Released: int(row[2].(int64)),
// 	// 		},
// 	// 	}
// 	// }

// 	// err = json.NewEncoder(response).Encode(results)
// 	// if err != nil {
// 	// 	log.Println("error writing search response:", err)
// 	// 	response.WriteHeader(500)
// 	// 	response.Write([]byte("An error occurred writing response"))
// 	// }
// }

// function movieHandler(w http.ResponseWriter, req *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	query := req.URL.Path[len("/movie/"):]
// 	cypher := `
// 	MATCH
// 		(movie:Movie {title:{title}})
// 	OPTIONAL MATCH
// 		(movie)<-[r]-(person:Person)
// 	WITH
// 		movie.title as title,
// 		collect({name:person.name, job:head(split(lower(type(r)),'_')), role:r.roles}) as cast
// 	LIMIT 1
// 	UNWIND cast as c
// 	RETURN title, c.name as name, c.job as job, c.role as role`

// 	db, err := driver.NewDriver().OpenNeo(neo4jURL)
// 	if err != nil {
// 		log.Println("error connecting to neo4j:", err)
// 		w.WriteHeader(500)
// 		w.Write([]byte("An error occurred connecting to the DB"))
// 		return
// 	}
// 	defer db.Close()

// 	data, _, _, err := db.QueryNeoAll(cypher, map[string]interface{}{"title": query})
// 	if err != nil {
// 		log.Println("error querying movie:", err)
// 		w.WriteHeader(500)
// 		w.Write([]byte("An error occurred querying the DB"))
// 		return
// 	} else if len(data) == 0 {
// 		w.WriteHeader(404)
// 		return
// 	}

// 	movie := Movie{
// 		Title: data[0][0].(string),
// 		Cast:  make([]Person, len(data)),
// 	}

// 	for idx, row := range data {
// 		movie.Cast[idx] = Person{
// 			Name: row[1].(string),
// 			Job:  row[2].(string),
// 		}
// 		if row[3] != nil {
// 			movie.Cast[idx].Role = interfaceSliceToString(row[3].([]interface{}))
// 		}
// 	}

// 	err = json.NewEncoder(w).Encode(movie)
// 	if err != nil {
// 		log.Println("error writing movie response:", err)
// 		w.WriteHeader(500)
// 		w.Write([]byte("An error occurred writing response"))
// 	}
// }

// function graphHandler(w http.ResponseWriter, req *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	limits := req.URL.Query()["limit"]
// 	limit := 50
// 	var err error
// 	if len(limits) > 0 {
// 		limit, err = strconv.Atoi(limits[0])
// 		if err != nil {
// 			w.WriteHeader(400)
// 			w.Write([]byte("Limit must be an integer"))
// 		}
// 	}

// 	cypher := `
// 	MATCH
// 		(m:Movie)<-[:ACTED_IN]-(a:Person)
// 	RETURN
// 		m.title as movie, collect(a.name) as cast
// 	LIMIT
// 		{limit}`

// 	db, err := driver.NewDriver().OpenNeo(neo4jURL)
// 	if err != nil {
// 		log.Println("error connecting to neo4j:", err)
// 		w.WriteHeader(500)
// 		w.Write([]byte("An error occurred connecting to the DB"))
// 		return
// 	}
// 	defer db.Close()

// 	stmt, err := db.PrepareNeo(cypher)
// 	if err != nil {
// 		log.Println("error preparing graph:", err)
// 		w.WriteHeader(500)
// 		w.Write([]byte("An error occurred querying the DB"))
// 		return
// 	}
// 	defer stmt.Close()

// 	rows, err := stmt.QueryNeo(map[string]interface{}{"limit": limit})
// 	if err != nil {
// 		log.Println("error querying graph:", err)
// 		w.WriteHeader(500)
// 		w.Write([]byte("An error occurred querying the DB"))
// 		return
// 	}

// 	d3Resp := D3Response{}
// 	row, _, err := rows.NextNeo()
// 	for row != nil && err == nil {
// 		title := row[0].(string)
// 		actors := interfaceSliceToString(row[1].([]interface{}))
// 		d3Resp.Nodes = append(d3Resp.Nodes, Node{Title: title, Label: "movie"})
// 		movIdx := len(d3Resp.Nodes) - 1
// 		for _, actor := range actors {
// 			idx := -1
// 			for i, node := range d3Resp.Nodes {
// 				if actor == node.Title && node.Label == "actor" {
// 					idx = i
// 					break
// 				}
// 			}
// 			if idx == -1 {
// 				d3Resp.Nodes = append(d3Resp.Nodes, Node{Title: actor, Label: "actor"})
// 				d3Resp.Links = append(d3Resp.Links, Link{Source: len(d3Resp.Nodes) - 1, Target: movIdx})
// 			} else {
// 				d3Resp.Links = append(d3Resp.Links, Link{Source: idx, Target: movIdx})
// 			}
// 		}
// 		row, _, err = rows.NextNeo()
// 	}

// 	if err != nil && err != io.EOF {
// 		log.Println("error querying graph:", err)
// 		w.WriteHeader(500)
// 		w.Write([]byte("An error occurred querying the DB"))
// 		return
// 	} else if len(d3Resp.Nodes) == 0 {
// 		w.WriteHeader(404)
// 		return
// 	}

// 	err = json.NewEncoder(w).Encode(d3Resp)
// 	if err != nil {
// 		log.Println("error writing graph response:", err)
// 		w.WriteHeader(500)
// 		w.Write([]byte("An error occurred writing response"))
// 	}
// }

// function init() {
// 	if os.Getenv("NEO4J_URL") != "" {
// 		neo4jURL = os.Getenv("NEO4J_URL")
// 		log.Printf("neo4j_URL = " + neo4jURL)
// 	}
// }

// function main() {
// 	serveMux := http.NewServeMux()
// 	serveMux.HandleFunc("/", defaultHandler)
// 	serveMux.HandleFunc("/search", searchHandler)
// 	serveMux.HandleFunc("/movie/", movieHandler)
// 	serveMux.HandleFunc("/graph", graphHandler)

// 	port := os.Getenv("PORT")
// 	if port == "" {
// 		port = "8080"
// 	}

// 	log.Printf("Starting server on port %s with neo4j %s", port, neo4jURL)
// 	panic(http.ListenAndServe(":"+port, serveMux))
// }