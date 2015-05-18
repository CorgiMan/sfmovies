# San Francisco Movies Challange

## Challange Description
Create a service that shows on a map where movies have been filmed in San Francisco. The user should be able to filter the view using autocompletion search.

## The San Francisco Movies API
I have chosen for the technical track and only provide a very basic front end at [](http://corgiman.infty.nl:12080). I have created a RESTful JSON API service that handles autocomplete, search and location based requests. The API (and its description) is located at [](http://corgiman.infty.nl/) The service uses [Google's Geoencoding API](https://developers.google.com/maps/documentation/geocoding/) for translating location names in gps coordinates and [The Open Movie Database API](http://www.omdbapi.com/) for movie info.

Goals: We want our service to be always available, reliable and blazingly fast. 

All the code included in this repository excecpt for some files in the front end are written by me. I did not write the the jquery files and `frontend/gmaps.js`.

### API Description
The API handles autocomple, search and location based request. The API can also be queried for specific movies or scenes by providing the IMDB movie id, or scene id.

- [](corgiman.infty.nl/status) returns the status of the api server that handled the request
- [](corgiman.infty.nl/movies/imdb_id/tt0028216) returns the movie info of the specified imdb id
- [](corgiman.infty.nl/scenes/scene_id/XXX) returns info of the specified scene id
- [](corgiman.infty.nl/complete?term=franc) auto complete the term.
- [](corgiman.infty.nl/search?term=francisco) searches for movie titles, film locations, release year, directors, production companies, distributers, writers and actors",
- [](http://corgiman.infty.nl/near?lat=37.76&lng=-122.39): "Search for film locations near the presented gps coordinates"

Use the callback parameter (?callback=XXX) on any request to return jsonp in stead of just json

## System Design
There are a couple of possible directions our application could scale in: the size of the source table, frequent updates of the data and the amount of requests per second. 

The data set that we work with is very small! The source table for our service only contains 1151 rows and is unlikely to grow very fast. The data set is also not likely to be updated very frequently. As a matter of fact, it hasn't been updated in years. Let's pretend however that updated when a new film is recorded in San Francisco.It seems like the amount of requests per second is the only probable direction this scenary could scale in so let's focus on that. 

We can deploy a load balancer to distribute the potentially high amounts of traffic amongst a set of nodes that handle API requests. Since the data set is so small we can easily fit dat into the memory of every web process of every API server. The API servers fetch the latest data from a MongoDB database which is updated once a day by a manager node. The manager node also monitors the API servers and restarts them if they fail to respond.

High availabillity and reliabillity: If one or a few of the API servers crash the load can be easily handled by the rest of them. The crashed servers are restarted by a monitoring script running on the manager node. If the database fails we can still serve data because the API servers store the data in memory. We have enough time to restart the database. 

Performance: Because the API servers store their data in memory they are extremely fast. 

With this design we can see that our goals of performance, high availabillity and reliabillity are met. The rest of this chapter discusses details of the implementation.

### Docker
Since I don't possess multiple machines, Docker seemed like a good choice to simulate the components that are part of our architecture (i.e. load balancer, mongodb, manager, API servers). I had no previous experience with Docker but because I heard a lot about the project, it seemed like an excellent oppertunity to learn. 

Docker only works on 64-bit machines and my linux system is only a 32-bit machine. I do, however, have a 64-bit windows machine and by installing boot2docker virtual machine we can use the Docker anyway. There was a big problem with the boot2docker vm: docker volumes didn't work properly for me. I wasn't able to access anything by disc, therefor, when building the images, I chose to download all necessary files with git. This is something that you probably avoid in the real world. The volume accessing issue also affected logging and cronjobs which are therefor completely managed within the containers.

### Load Balancing
Nginx seems like an excellent choice to be a load balancer. The Nginx Dockerfile hosted on dockerhub didn't work for my so I made a custom one `sfmovies/docker/loadbalancer`. The container is run it as a daemon and keeps bash on the foreground open so that I can reload config file changes without being disconnected for a moment. I found out that Nginx also reloads the config file when send a "HUP" signal with killall. We can use this to let Nginx run in the forground of the container. The problem I have with the volume access, however, still disables me from running nginx in the forground of the container because I don't have access to the config file from outside the container.

The ip addresses and ports of the API servers need to be added to the sfmovies/nginx.conf file so that nginx can load balance between them.

### Mongodb
The mongodb image hosted on dockerhub fulfills our service's needs. For security we want to add `bind_ip = 127.0.0.1` to the mongodb.conf so that mongo only listens to localhost.
 
### API Server
The API server just installs the `gocode/APIserver` with `go install`. Implementations details are described in a later chapter.

The host of the API server containers needs to install the public ssh key of the manager server. The manager server uses ssh to monitor and restart API servers if necassery. On start or restart the API server gets the latest dataset from the mongodb database.

### Manager
The manager updates the dataset every day in the night and monitors all the API servers every minute with cronjobs. After the dataset is updated, a rolling restart of all the API servers is performed by a shell script `docker/manager/rolling_restart.sh`. The monitoring is also done by `docker/manager/monitor.sh`. A single restart is performed by `docker/manager/restart.sh`. To manage the API server containers the manager needs a list of ip:port pairs nodes: `docker/manager/API_servers`.

## API Server Implementation (Go)
On initialization, the program fetches the API data from Mongodb. The program uses go's in build web server to handle the requests. For the search and autocomplete requests we use a trie. We could use a quad tree for the location based searches, but with only ~1200 points-of-interrest we don't gain much from a quad tree approach. I have chosen for simplicity instead of a very small gain. All the handlers are wrapped in a callback handler that serves JSONP if the ?callback parameter is set. JSONP is necessary when requesting data from a different domain than the domain that hosts the front end.


## Front End
I programmed a small front end for testing and as a prove of concept. Although the front end is hosted in the same domain as the API servers, you can easily check that it works outside this domain by downloading the front end and opening the index.html with your browser. It uses jQuery, the Google Maps API and gmaps.js.

Just a day before the deadline I started hosting the front end on corgiman.infty.nl:12080. There seemed to be a problem with the movies posters that were loaded from IMDB. Where running on localhost on my machine gave no problems, the hosted version gave 403 frobidden access status codes. As it turns out IMDB doesn't like third parties to use their images. This is a flaw the API and I had to solve quickly. I started comparing requests headers from both my local machine and the hosted machine. I dit this by setting up a reverse proxy which send the movie poster requests to port 12081 and I listened to the requests using `nc -l 12081`. I found out that the main difference between the requests was the appearance of the Referer header. I set up a reversed proxy with Nginx to IMDB which I configured to remove the Referer header and this solved the problem.

Happy to see that this works just before the deadline of this challange I must place a few remarks. In production you would not want to use a reverse proxy because IMDB might block the server ip. We could cache the images with Nginx or store them in Mongodb. There are not a lot of movies recorded in San Fransisco so we could just download and store all the movie posters we need. Another solution is to use the OMDB poster API for which you need to donate to get a key.


## Installation
If you want to test the system on a single follow the commands below. The configuration file docker/loadbalancer/nginx.conf specifies ports 12001, 12002 and 12003 for load balancing. If you want to run more API server nodes you need to attach to the nginx container, add the adress of the new API server to the /home/nginx.conf and reload the config file with the nginx -s reload command.

    git clone https://github.com/CorgiMan/sfmovies.git
    
    docker build -t sfmovies/nginx ./sfmovies/docker/loadbalancer
    docker build -t sfmovies/APIserver ./sfmovies/docker/APIserver
    docker build -t sfmovies/manager ./sfmovies/docker/manager
    
    docker run -dit -p 80:80 sfmovies/nginx

    docker run -d -p 27017:27017 -name mongodb mongo

    docker run -d -p 12100:80 sfmovies/manager

    docker run -d -p 12001:80 sfmovies/APIserver
    docker run -d -p 12002:80 sfmovies/APIserver
    docker run -d -p 12003:80 sfmovies/APIserver




