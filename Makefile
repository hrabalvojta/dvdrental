prep:
	docker stop postgres; \
	sleep 1; \
	docker run --rm --name postgres -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=dvdrental -p5432:5432 -d postgres:alpine

sleep:
	sleep 3

psql:
	PGPASSWORD=secret psql -hlocalhost -p5432 -Upostgres -d dvdrental

psql_reset: prep sleep psql

update:
	go get -u all; \
	go mod tidy -compat=1.17
	
compile:
	echo "Compiling for every OS and Platform"
	for arch in "386" "amd64" "arm64"; do \
			GOOS=linux GOARCH=$$arch go build -o bin/films-linux-$$arch cmd/films/main.go; \
			GOOS=windows GOARCH=$$arch go build -o bin/films-windows-$$arch.exe cmd/films/main.go; \
	done

test:
	for i in `seq 1 5`; do \
		curl -X POST -d'{"a":10,"b":10}' localhost:8081/sum; \
		curl -X POST -d'{"a":"10","b":"10"}' localhost:8081/concat; \
		sleep 1; \
	done; \
	curl localhost:8080/metrics | grep -i dvdrental

