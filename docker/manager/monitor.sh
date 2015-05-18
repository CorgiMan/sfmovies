#!/bin/sh
# monitors the api servers listed in the api_servers file
# for every api server this script requests the status of the server
# if we receive a statuscode of 200 we continue to the next api server
# in other cases we try and restart the server using the restart.sh script

echo start monitoring...

all_running=1
cat api_servers | while read address; do
    address=$(echo $address | awk '{print $1}')
    
    if [ $(curl -s -o /dev/null -w "%{http_code}" http://$address/status) != 200 ]; then
      echo no response from $address. Trying restart...
      /bin/sh restart.sh $address
      all_running=0
    fi
done

if [ $all_running == 1 ]; then
  echo all servers are running
fi
