FROM ubuntu

RUN \
  apt-get update && \
  apt-get install -y software-properties-common python-software-properties && \
  add-apt-repository -y ppa:nginx/stable && \
  apt-get install -y nginx && \
  apt-get install -y nano && \
  apt-get isntall -y curl && \
  curl https://raw.githubusercontent.com/CorgiMan/sfmovies/master/loadbalancer/nginx.conf > /home/nginx.conf && \
  rm -rf /var/lib/apt/lists/* && \
  #echo "\ndaemon off;" >> /etc/nginx/nginx.conf && \
  chown -R www-data:www-data /var/lib/nginx 

VOLUME ["/etc/nginx/sites-enabled", "/etc/nginx/certs", "/etc/nginx/conf.d", "/var/log/nginx", "/var/www/html"]


# run nginx as a daemon and run bash in forground so that we can
# change and reload the config file without losing trafic

CMD nginx -c /home/nginx.conf && \
    /bin/bash
    
EXPOSE 80


# docker build -t sfmovies/nginx $GOPATH/src/github.com/sfmovies/loadbalancer
# docker run -dit -p 80:80 sfmovies/nginx

# docker run -dit -p 80:80 -v //c/Users/Wieger/Desktop/sfmovies/loadbalancer:/home/ sfmovies/nginx
