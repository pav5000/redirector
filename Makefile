
IMAGE:=redirector

build:
	sudo docker build -t redirector .

run:
	go run github.com/pav5000/redirector/cmd/redirector

start:
	sudo docker run -d \
		--name=${IMAGE} \
		--restart=always \
		--network=host \
		-v `pwd`/config.yml:/config.yml \
		${IMAGE}

stop:
	sudo docker stop ${IMAGE}; true
	sudo docker rm -f ${IMAGE}; true

restart: stop start
	echo "restarted"

logs:
	sudo docker logs --tail 100 -f ${IMAGE}

status:
	sudo docker stats ${IMAGE} || echo "Status: stopped"
