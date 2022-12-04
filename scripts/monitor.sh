# 监控硬盘状态，/dev/vdb1的Avail值小于80G时, kill掉所有download.sh进程
# 用法：在crontab中添加一条定时任务，每分钟执行一次
# 例如：* * * * * root /bin/bash /root/ishare/scripts/monitor.sh>>/home/ishare/scripts/monitor.log
# 重启crontab服务：service crond restart

# 0. 打印当前时间
limit=70
echo "----------------------------------------"
echo "Current time: $(date)"
# 1. 获取/dev/vdb1的Avail值
avail=$(df -h | grep /dev/vdb1 | awk '{print $4}' | sed 's/G//g')
echo "avail: $avail"
# 2. 判断Avail值是否小于80
if [ $avail -lt $limit ]; then
    echo "avail is less than $limit"
    # 3. 获取所有download.sh进程的pid
    pids=$(ps -ef | grep download.sh | grep -v grep | awk '{print $2}')
    echo "find downdload pids: $pids"
    # 4. kill掉所有download.sh进程
    for pid in $pids; do
        echo "kill $pid"
        kill -9 $pid
    done
else
    echo "avail is greater than $limit"
fi