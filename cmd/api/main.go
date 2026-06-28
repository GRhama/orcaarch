package main

import (
	"fmt"
	"net/http"

	"orcaarch/internal/api"
	"orcaarch/internal/service"
)

func main() {
	result, err := service.Process(1000)
	if err != nil {
		panic(err)
	}
	fmt.Println("OrcaArch API — logistics scenario, 1000 records")
	fmt.Println()
	fmt.Println("  GET http://localhost:8080/api/v1/reports/inventory")
	fmt.Println("  GET http://localhost:8080/api/v1/reports/ledger")
	fmt.Println("  GET http://localhost:8080/api/v1/reports/risk")
	fmt.Println()
	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", api.NewServer(result))
}
