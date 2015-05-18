#!/bin/sh
# restart a container remotely. input format: 123.45.56.78:12345
# if the container does not exist we try to spawn a new container that listens to the correct port

host=$(echo $1 |cut -d':' -f1)
port=$(echo $1 |cut -d':' -f2)

# get the container id that belongs to the host and port
cid=$(0</dev/null ssh $host docker ps | grep sfmovies/apiserver | grep :$port | awk '{print $1}')
if [ "$cid" == "" ]; then
        echo "no container found at $host:$port"
        echo $(0</dev/null ssh $host docker run -d -p $port:80 sfmovies/apiserver) started
else
        echo $(0</dev/null ssh $host docker restart $cid) restarted
fi

sleep 1s
while [ $(curl -s -o /dev/null -w "%{http_code}" http://$host:$port/status) != 200 ]; do
        echo trying to connect to new server...
        sleep 1s
done

echo server status:
curl http://$host:$port/status
echo ""
