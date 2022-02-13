prep:
	docker run --rm --name postgres -e POSTGRES_PASSWORD=secret -e POSTGRES_DB=dvdrental -p5432:5432 -d postgres:alpine

compile:
	echo "Compiling for every OS and Platform"
	for arch in "386" "amd64"; do \
			GOOS=linux GOARCH=$$arch go build -o bin/films-linux-$$arch cmd/films/main.go; \
			GOOS=windows GOARCH=$$arch go build -o bin/films-windows-$$arch.exe cmd/films/main.go; \
	done
