docker run -p 80:80 -p 3336:3306 \
 -v /data/lamp/www:/var/www \
 -v /data/lamp/apache-conf/apache2.conf:/etc/apache2/apache2.conf \
 -v /data/lamp/mysql-conf/my.cnf:/etc/mysql/my.cnf \
 -t -i linode/lamp /bin/bash
