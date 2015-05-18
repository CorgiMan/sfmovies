# San Francisco Movies Challenge

## Challenge Description
[Uber San Francisco Movies Challenge](https://github.com/uber/coding-challenge-tools/blob/master/coding_challenge.md):
"Create a service that shows on a map where movies have been filmed in San Francisco. The user should be able to filter the view using auto-completion search."

[Data Source](https://data.sfgov.org/Arts-Culture-and-Recreation-/Film-Locations-in-San-Francisco/yitu-d5am)

## [The San Francisco Movies API](http://corgiman.infty.nl/)
I have chosen for the technical track and only provide a very basic [**front end**](http://corgiman.infty.nl:12080) at [corgiman.infty.nl:12080](http://corgiman.infty.nl:12080). I have created a RESTful JSON API service that handles auto-complete, search and location based requests. The API is located at [corgiman.infty.nl](http://corgiman.infty.nl/). The service uses [Google's Geo-encoding API](https://developers.google.com/maps/documentation/geocoding/) for translating location names in gps coordinates and [The Open Movie Database API](http://www.omdbapi.com/) for movie info.

Goals: We want our service to be always available, reliable and blazingly fast. 

All the code included in this repository except for some files in the front end are written by me. I did not write `frontend/gmaps.js` and the jQuery files.

### API Description
The API handles auto-complete, search and location based request. The API can also be queried for specific movies by providing the IMDB movie id.

- [corgiman.infty.nl/status](http://corgiman.infty.nl/status) The status of the API server that handled the request
- [corgiman.infty.nl/movies/tt0028216](http://corgiman.infty.nl/movies/imdb_id/tt0028216) Movie info of the specified IMDB id
- [corgiman.infty.nl/complete?term=franc](http://corgiman.infty.nl/complete?term=franc) Auto-complete the term parameter
- [corgiman.infty.nl/search?q=francisco](http://corgiman.infty.nl/search?term=francisco) Searches for movie titles, film locations, release year, directors, production companies, distributors, writers and actors
- [corgiman.infty.nl/near?lat=37.76&lng=-122.39](http://corgiman.infty.nl/near?lat=37.76&lng=-122.39) Search for film locations near the presented gps coordinates

Use the callback parameter (?callback=XXX) on any request to return JSONP instead of just JSON.

## System Design
There are a couple of possible directions our application could scale in: the size of the source table, frequent updates of the data and the amount of requests per second. 

The data set that we work with is very small! The source table for our service only contains 1151 rows and is unlikely to grow very fast. The data set is also not likely to be updated very frequently. As a matter of fact, it hasn't been updated in years. Let's pretend however that the source is updated when a new film is recorded in San Francisco. This rules out the first two scaling directions. It seems like the amount of requests per second is the most probable direction the service could scale in, so let's focus on that. 

We can deploy a load balancer to distribute the potentially high amounts of traffic amongst a set of nodes that handle API requests. Since the data set is so small, we can easily fit it into the memory of every web process of every API server. The API servers fetch the latest data from a MongoDB database which is updated once a day by a manager node. The manager node also monitors the API servers and restarts them if they fail to respond.

With this design our goals of performance, high availability and reliability are met.
- Performance: Because the API servers store their data in memory they are extremely fast. 
- High availability and reliability: If one or a few of the API servers crashes, the load can be easily handled by the rest of them. The crashed servers are restarted by a monitoring script running on the manager node. If the database fails we can still serve data because the API servers store the data in memory and only access the database on initialization.

The rest of this chapter discusses some details of the implementation.

### Docker
Since I don't possess multiple machines, Docker seems like a good choice to build the components of our architecture (i.e. load balancer, MongoDB, manager, API servers). I had no previous experience with Docker but because I heard a lot about the project, it seemed like an excellent learning opportunity. 

Docker only works on 64-bit machines and my linux system is only a 32-bit machine. I do, however, have a 64-bit windows machine and by installing the boot2docker virtual machine we can use the Docker anyway. There was a big problem with the boot2docker vm: docker volumes didn't work properly for me. I wasn't able to access anything by disc, therefore, when building the images, I chose to download all necessary files from this repo with git. The volume accessing issue also affected logging and cronjobs which are therefore completely managed within the containers.

### Load Balancing
Nginx is an excellent load balancer. The Nginx Dockerfile hosted on Docker Hub didn't work for my so I made a custom one `sfmovies/docker/loadbalancer`. The container is run as a daemon and keeps bash on the foreground so that I can reload config file changes without being disconnected for a moment. I found out that Nginx also reloads the config file when send a "HUP" signal with killall. We could use this to let Nginx run in the foreground of the container. The problem I have with the volume access, however, still disables me from running Nginx in the foreground of the container because I don't have access to the config file from outside the container.

The ip addresses and ports of the API servers need to be added to the sfmovies/nginx.conf file so that Nginx can load balance between them.

### MongoDB
The MongoDB image hosted on Docker Hub fulfills our service's needs.
 
### API Server
The API server just installs the `gocode/apiserver` with `go install`. Implementation details are given in a later chapter. On initialization, the API server fetches the latest dataset from the MongoDB database.

The host that runs the API server containers needs to install the public ssh key of the manager server. The manager server uses ssh to monitor and restart API servers if necessary. 

### Manager
The manager updates the dataset every night at 4 o' clock and monitors the API servers every minute with the use of cronjobs. After the dataset is updated, a rolling restart of all the API servers is performed by a shell script `docker/manager/rolling_restart.sh`. The monitoring is done by `docker/manager/monitor.sh`. These two scripts use a script for the restart of a single container `docker/manager/restart.sh`. To manage the nodes, the manager needs a list of ip:port pairs to the API server containers: `docker/manager/API_servers`.

## API Server Implementation (Go)
On initialization, the program fetches the latest API data from MongoDB. The program uses go's build-in web server to handle the requests. For the search and auto-complete requests we use a trie. We could use a quad tree for the location based searches, but with only ~1200 points-of-interest we don't gain much from a quad tree approach. I have chosen for simplicity instead of a very small gain. All the handlers are wrapped in a callback handler that serves JSONP if the `?callback` parameter is set. JSONP is used when requesting data from a different domain than the domain that hosts the front end.


## Front End
I programmed a small front end for testing and as a prove of concept. Although the front end is hosted in the same domain as the API servers, you can easily check that it works outside this domain by downloading the front end and opening the index.html with your browser. It uses jQuery, the Google Maps API and gmaps.js.

Just a day before the deadline I started hosting the front end on corgiman.infty.nl:12080. There was a problem with the movie posters that were loaded from IMDB. Where running on localhost on my machine gave no problems, the hosted version gave 403 forbidden access status codes. As it turns out, IMDB doesn't like third parties to use their images. This is a flaw in the API and one I had to solve quickly. I started comparing requests headers from both my local machine and the hosted machine. I did this by setting up a reverse proxy with Nginx which sends the movie poster requests to port 12081. I listened to the requests on this port using `nc -l 12081`. I found out that the main difference between the requests was the appearance of the Referer header. I set up a reversed proxy with Nginx to IMDB which I configured to remove the Referer header and this solved the problem.

Happy to see that this works just before the deadline of this challenge I must place a few remarks. In production you would not want to use a reverse proxy because IMDB might block the server ip. We could cache the images with Nginx or store them in MongoDB. There are not a lot of movies recorded in San Fransisco so we could just download and store all the movie posters we need. Another solution is to use the OMDB poster API for which you need to donate to get a key.


## Installation
If you want to test the system on a single machine follow the commands below. The configuration file docker/loadbalancer/nginx.conf specifies ports 12001, 12002 and 12003 for load balancing. If you want to run more API server nodes you need to attach to the Nginx container, add the address of the new API server to the /home/nginx.conf and reload the config file with the Nginx -s reload command. You also need to also add them to /manager/api_servers. I limited the amount of rows the service scans from the source table, because the Google Geo-encoding API limits the amount of request to 2500 per day. This is enough for our needs but if everybody installs the system we might run into some problems.

    git clone https://github.com/CorgiMan/sfmovies.git
    
    docker build -t sfmovies/nginx ./sfmovies/docker/loadbalancer
    docker build -t sfmovies/APIserver ./sfmovies/docker/APIserver
    docker build -t sfmovies/manager ./sfmovies/docker/manager
    
    docker run -dit -p 80:80 sfmovies/nginx

    docker run -d -p 27017:27017 -name MongoDB mongo

    docker run -d -p 12100:80 sfmovies/manager

    docker run -d -p 12001:80 sfmovies/APIserver
    docker run -d -p 12002:80 sfmovies/APIserver
    docker run -d -p 12003:80 sfmovies/APIserver

### Improvements
- Don't use windows' book2docker. Use linux instead so that we can make use of the volumes. I faced a lot of nasty problems with the boot2docker setup, but it was the only option that I had at the moment.
- Finish processes gracefully. After the daily database update we perform a rolling restart. This could lead to some unfinished requests. With go, we can interrupt the SIGTERM signal that is send when docker restarts the containers, and gracefully finish the requests first.
- As described in the Front End section, it is a bad idea to reverse proxy the IMDB movie posters. 
- Search is not very versatile yet. An easy improvement to search would be to allow multiple words in the query and return the intersection of the individual results.
- Auto-complete responds only with words. I'd like to change it so that if you search for "adam", the API auto-completes it to "Adam Sandler". To implement this the TrieNode should store a list of strings under every node. It takes some extra work because there are often multiple actors in a single string. e.g. "Adam Sandler, Drew Barrymore, Rob Schneider, Sean Astin" should be split into: ["Adam Sandler", "Drew Barrymore", "Rob Schneider", "Sean Astin"]

