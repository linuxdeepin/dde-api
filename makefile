run:=go build -o .out && rm .out

all: 
	cd graphic/ && $(run)
	cd logger/ && $(run)
	cd mousearea/ && $(run)
	cd pinyin-search/ && $(run)
	cd set-date-time/ && $(run)


update:
	sudo apt-get update && sudo apt-get install dde-go-dbus-factory go-dlib
