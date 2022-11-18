本部分脚本用于下载openAlex的数据

- url_compute.py 给定下载数据的大小download_size，计算可以下载的url存入urls中

- download.sh 使用urls中的url下载数据，下载的数据存入第一个参数标明的路径中，日志打入download.log中