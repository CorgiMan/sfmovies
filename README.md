# sfmovies 


## San Francisco Movies API
api description ...
link to backend infty.nl:12000
link to frontend infty.nl:12200

## Used Services
google geocoding
omdb

## System Design
I used docker ...

### loadbalancer
nginx dockerfile hosted on dockerhub didn't work for my so I made one. I run it as a daemon and keep bash open so that I can reload config file changes without being disconnected for a moment.

### api server
usually lifetime of 24 hours ...

### mongodb
docker run -d -p 27017:27017 -name mongodb mongo

### monitor / data update server


## API Server Implementation
stores data in memory
uses trie for search and auto complete.

### Trie (search and auto complete queries)

### Quad Tree (near queries)
Could use a quad tree for the search, but with only ~1200 points-of-interrest we don't gain much from a quad tree approach. I have chosen for simplicity in stead of a very small gain.

## Front End
Small front end for testing and as a prove of concept. The front end uses the ?callback parameter so that the api server will return jsonp. This is necessary when requesting data from a different domain than the domain that hosts the front end. Although the front end is hosted in the same domain as the api servers, you can easily check that it works outside this domain by downloading the front end and opening the index.html with your browser.

