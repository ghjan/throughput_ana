docker run --name redis -p 6379:6379 -v /data/throughput/logs/redis/:/var/log/redis/ -v /data/throughput/work_run/redis/:/var/run/redis/ -v /data/throughput/work_run/redis_db:/var/lib/redis -v /data/gopath/src/github.com/ghjan/throughput_ana/redis/redis.conf:/etc/redis/redis.conf -dit redis:3.2
