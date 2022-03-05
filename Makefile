build_raspi:
	rm -rf build/arm/*
	docker run --rm \
		-v "$(shell pwd):/app" \
		-v "/tmp/go-docker:/go" \
		-w /app \
		--platform=linux/arm64 golang:1.17.6 \
		sh -c "apt update && apt install -y libpcap-dev; ls -lah;\
		env GOOS=linux; \
			CGO_ENABLED=1; \
			GOARCH=arm64; \
			GOARM=7; \
			go build -o build/arm/godhcpdump main.go"