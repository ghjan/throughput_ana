#lamp_init.sh
docker run -p 80:80 -p 3336:3306 -v /data/lamp/www:/var/www -t -i linode/lamp /bin/bash