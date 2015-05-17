#!/bin/sh

cat api_servers | while read address; do
    echo restarting container at $address
    /bin/sh restart.sh $address
    echo ""
done
echo "succesfully restarted servers"
