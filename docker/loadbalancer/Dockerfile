FROM ubuntu

RUN \
  apt-get update && \
  apt-get install -y software-properties-common python-software-properties && \
  add-apt-repository -y ppa:nginx/stable && \
  apt-get install -y nginx && \
  apt-get install -y nano && \
  apt-get install -y curl && \
  curl https://raw.githubusercontent.com/CorgiMan/sfmovies/master/docker/loadbalancer/nginx.conf > /home/nginx.conf && \
  rm -rf /var/lib/apt/lists/* && \
  #echo "\ndaemon off;" >> /etc/nginx/nginx.conf && \
  chown -R www-data:www-data /var/lib/nginx 

VOLUME ["/etc/nginx/sites-enabled", "/etc/nginx/certs", "/etc/nginx/conf.d", "/var/log/nginx", "/var/www/html"]


# run nginx as a daemon and run bash in forground so that we can
# change and reload the config file without losing trafic

CMD nginx -c /home/nginx.conf && \
    /bin/bash
    
EXPOSE 80
