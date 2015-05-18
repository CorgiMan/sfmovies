# San Francisco Movies Challange

## Challange Description

## Solution
Geo encoding. omdb imdb
Goals: We want our service to be always available, reliable and blazingly fast. 

## San Francisco Movies API
api description ...
link to backend infty.nl:12000
link to frontend infty.nl:12200

## Used Services
google geocoding
omdb
imdb

## System Design
There are a couple of possible directions our application could scale in: the size of the source table, frequent updates of the data and the amount of requests per second. 

The data set that we work with is very small! The source table for our service only contains 1151 rows and is unlikely to grow very fast. The data set is also not likely to be updated very frequently. As a matter of fact, it hasn't been updated in years. Let's pretend however that updated when a new film is recorded in San Francisco.It seems like the amount of requests per second is the only probable direction this scenary could scale in so let's focus on that. 

We can deploy a load balancer to distribute the potentially high amounts of traffic amongst a set of nodes that handle api requests. Since the data set is so small we can easily fit dat into the memory of every web process of every api server. The api servers fetch the latest data from a MongoDB database which is updated once a day by a manager node. The manager node also monitors the api servers and restarts them if they fail to respond.

High availabillity and reliabillity: If one or a few of the api servers crash the load can be easily handled by the rest of them. The crashed servers are restarted by a monitoring script running on the manager node. If the database fails we can still serve data because the api servers store the data in memory. We have enough time to restart the database. 

Performance: Because the api servers store their data in memory they are extremely fast. 

With this design we can see that our goals of performance, high availabillity and reliabillity are met. The rest of this chapter discusses details of the implementation.

### Docker
Since I don't possess multiple machines, Docker seemed like a good choice to simulate the components that are part of our architecture (i.e. load balancer, mongodb, manager, api servers). I had no previous experience with Docker but because I heard a lot about the project, it seemed like an excellent oppertunity to learn. 

Docker only works on 64-bit machines and my linux system is only a 32-bit machine. I do, however, have a 64-bit windows machine and by installing boot2docker virtual machine we can use the Docker anyway. There was a big problem with the boot2docker vm: docker volumes didn't work properly for me. I wasn't able to access anything by disc, therefor, when building the images, I chose to download all necessary files with git. This is something that you probably avoid in the real world. The volume accessing issue also affected logging and cronjobs which are therefor completely managed within the containers.

### Load Balancing
Nginx seems like an excellent choice to be a load balancer. The Nginx Dockerfile hosted on dockerhub didn't work for my so I made a custom one `sfmovies/docker/loadbalancer`. The container is run it as a daemon and keeps bash on the foreground open so that I can reload config file changes without being disconnected for a moment. I found out that Nginx also reloads the config file when send a "HUP" signal with killall. We can use this to let Nginx run in the forground of the container. The problem I have with the volume access, however, still disables me from running nginx in the forground of the container because I don't have access to the config file from outside the container.

The ip addresses and ports of the api servers need to be added to the sfmovies/nginx.conf file so that nginx can load balance between them.

### Mongodb
The mongodb image hosted on dockerhub fulfills our service's needs. For security we want to add `bind_ip = 127.0.0.1` to the mongodb.conf so that mongo only listens to localhost.
 
### API Server
The api server just installs the `gocode/apiserver` with `go install`. Implementations details are described in a later chapter.

The host of the api server containers needs to install the public ssh key of the manager server. The manager server uses ssh to monitor and restart api servers if necassery. On start or restart the api server gets the latest dataset from the mongodb database.

### Manager
The manager updates the dataset every day in the night and monitors all the api servers every minute with cronjobs. After the dataset is updated, a rolling restart of all the api servers is performed by a shell script `docker/manager/rolling_restart.sh`. The monitoring is also done by `docker/manager/monitor.sh`. A single restart is performed by `docker/manager/restart.sh`. To manage the api server containers the manager needs a list of ip:port pairs nodes: `docker/manager/api_servers`.

## API Server Implementation (Go)
stores data in memory
uses trie for search and auto complete.

### Trie (search and auto complete queries)

### Quad Tree (near queries)
Could use a quad tree for the search, but with only ~1200 points-of-interrest we don't gain much from a quad tree approach. I have chosen for simplicity in stead of a very small gain.

## Front End
Small front end for testing and as a prove of concept. The front end uses the ?callback parameter so that the api server will return jsonp. This is necessary when requesting data from a different domain than the domain that hosts the front end. Although the front end is hosted in the same domain as the api servers, you can easily check that it works outside this domain by downloading the front end and opening the index.html with your browser.



### installation
If you want to test the system on a single follow the commands below. The configuration file docker/loadbalancer/nginx.conf specifies ports 12001, 12002 and 12003 for load balancing. If you want to run more api server nodes you need to attach to the nginx container, add the adress of the new api server to the /home/nginx.conf and reload the config file with the nginx -s reload command.

    git clone https://github.com/CorgiMan/sfmovies.git
    
    ip=$(ifconfig eth0 | grep 'inet addr:' | cut -d: -f2 | awk '{ print $1})

    docker build -t sfmovies/nginx ./sfmovies/docker/loadbalancer
    docker build -t sfmovies/apiserver ./sfmovies/docker/apiserver
    docker build -t sfmovies/manager ./sfmovies/docker/manager
    
    docker run -dit -p 80:80 sfmovies/nginx

    docker run -d -p 27017:27017 -name mongodb mongo

    docker run -d -p 12100:80 sfmovies/manager

    docker run -d -p 12001:80 sfmovies/apiserver
    docker run -d -p 12002:80 sfmovies/apiserver
    docker run -d -p 12003:80 sfmovies/apiserver

## War story
Just one day before the deadline I started hosting the front end on corgiman.infty.nl:12080. There seemed to be a problem with the movies posters that were loaded from imdb. Where running on localhost on my machine gave no problems, the hosted version gave 403 frobidden access status codes. As it turns out IMDB doesn't like third parties to use their images. This is a flaw the api and I had to solve quickly. I started comparing requests headers from both my local machine and the hosted machine. I dit this by setting up a reverse proxy which send the movie poster requests to port 12081 and I listened to the requests using `nc -l 12081`. As I found out that the main difference between the requests was the appearance of the Referer header. I set up a reversed proxy with Nginx which I configured to remove the Referer header and this solved the problem. 

Happy to see that this works just before the deadline of this challange I must place a few remarks. In production you would not want to use a reverse proxy because imdb might block the ip server. There are a few sollutions. Cache the images with Nginx or store them in Mongodb. There are not a lot of movies recorded in San Fransisco so we could just download and store all the movie posters we need. Another solution is to use the omdb poster api. I didn't have a key for the api and ran out of time at this point so I haven't implemented this.




