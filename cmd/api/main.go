package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"orcaarch/internal/api"
	"orcaarch/internal/service"
)

func main() {
	volume := flag.Int("volume", 1000, "number of records to process")
	flag.Parse()
	if *volume <= 0 {
		log.Fatal("volume must be greater than zero")
	}

	result, err := service.Process(*volume)
	if err != nil {
		panic(err)
	}
	fmt.Printf("OrcaArch API — logistics scenario, %d records\n", *volume)
	fmt.Println()
	fmt.Println("  GET http://localhost:8080/api/v1/reports/inventory")
	fmt.Println("  GET http://localhost:8080/api/v1/reports/ledger")
	fmt.Println("  GET http://localhost:8080/api/v1/reports/risk")
	fmt.Println()
	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", api.NewServer(result))
}
