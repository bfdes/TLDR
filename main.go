package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bradfitz/gomemcache/memcache"
	_ "github.com/lib/pq"
)

func getOrElse(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultValue
}

var (
	port       = getOrElse("PORT", "8080")
	cacheHost  = getOrElse("MEMCACHED_HOST", "localhost")
	cachePort  = getOrElse("MEMCACHED_PORT", "11211")
	dbHost     = getOrElse("POSTGRES_HOST", "localhost")
	dbPort     = getOrElse("POSTGRES_PORT", "5432")
	dbUser     = getOrElse("POSTGRES_USER", "postgres")
	dbPassword = getOrElse("POSTGRES_PASSWORD", "pass")
	dbName     = getOrElse("POSTGRES_DB", "url-shortener")
)

func initCache(host, port string) *memcache.Client {
	connStr := fmt.Sprintf("%s:%s", host, port)
	return memcache.New(connStr)
}

func initDb(host, port, user, password, dbName string) *sql.DB {
	connStr := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	query := `
		CREATE TABLE IF NOT EXISTS links(
			id BIGSERIAL PRIMARY KEY,
			url VARCHAR
		);
	`
	_, err = db.Exec(query)
	if err != nil {
		panic(err)
	}
	return db
}

func main() {
	cache := initCache(cacheHost, cachePort)
	db := initDb(dbHost, dbPort, dbUser, dbPassword, dbName)
	defer db.Close()
	linkService := linkService{cache, db}
	http.Handle("/api/links", CreateLinkHandler(linkService))
	http.Handle("/", RedirectHandler(linkService))
	port := fmt.Sprintf(":%s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
