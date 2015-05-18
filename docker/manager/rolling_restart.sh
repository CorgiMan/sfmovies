#!/bin/sh
# performs a rolling restart of all the api servers listed in the api_servers file
# the restart.sh script restarts a container and finishes 
# only when the restarted api server is handling requests

cat api_servers | while read address; do
    address=$(echo $address | awk '{print $1}')
    echo restarting container at $address
    /bin/sh restart.sh $address
    echo ""
done
echo "succesfully restarted servers"
