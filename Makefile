build:
	go build

run:
	go run main.go

test:
	go test ./...

make eval:
	time go run main.go --no-daemon true
	xsv sort -NRs weight meaningful-ngrams.csv -o meaningful-ngrams.csv
	csvlens meaningful-ngrams.csv
