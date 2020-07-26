run:
		go run main.go

build:
		go build -o bundle

start:
		pm2 start --name bundle ./bundle

stop:
		pm2 stop bundle
	
restart:
		pm2 restart bundle

logs:
		pm2 logs
install:
		go get -v