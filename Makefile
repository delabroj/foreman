.PHONY: buildpi

buildpi:
	@GOOS=linux GOARCH=arm GOARM=7 go build && zip install.zip foreman foreman.service install.sh && rm foreman
