docker run --name nginx -p 8000:8000 -v /data/nginx/logs:/etc/nginx/logs -v /data/nginx/nginx.conf:/etc/nginx/nginx.conf -dit nginx
