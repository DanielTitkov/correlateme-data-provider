run: 
	go run cmd/provider/main.go

gen:
	go run cmd/gen/main.go

run-backfill:
	go run cmd/provider/main.go -start=2023-09-01 -end=2023-09-30 -metrics="Random from 1 to 10|Another Metric"