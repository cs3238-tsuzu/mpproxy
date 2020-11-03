.PHONY: mpproxy
mpproxy:
	go build -o mpproxy ./cmd/mpproxy

upload:
	GOOS=linux go build -o mpproxy_linux ./cmd/mpproxy
	rsync -avh mpproxy_linux lightsail:mpproxy
