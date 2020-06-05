sudo docker-compose up -d
docker-compose exec rabbitmq  /bin/bash -c "apt-get update && apt-get -y install tzdata"
sudo docker-compose logs -f servicewa
