# -*- coding: utf-8 -*-
# 读取manifest文件夹下的json文件
import json
import os
import re

# 本脚本需要在scripts文件夹下运行,该文件夹内有若干相关文件
# download_size 是下载数据量大小参数
# openAlex 一共330G压缩数据，解压后1.6T
# 那么gzip的压缩率算作是5 我们的磁盘空间一共有950G，可用空间为888G，所以我们的下载数据量大小为800/5=160G(100G留作余量)
# 由于openAlex数据冗余的情况比较明显，可能经过数据清洗后，我们的数据量会更小。后期可以再进行调整
GB = 1024 * 1024 * 1024
ALL = 400
download_size = 102 * GB
manifests_path = os.path.join('.', 'manifests')
entries = []
urls = []
# 将所有json文件中的entries合并到一个list中
# for directory, subdir_list, file_list in os.walk(manifests_path):
#     for file in file_list:
with open(os.path.join(manifests_path, "works.json"), 'r') as f:
    data = json.load(f)
    entries.extend(data['entries'])
# 将entries按照url中updated_date=xxxx-xx-xx后的时间进行排序
date_pattern = r'\d{4}-\d{2}-\d{2}'
entries.sort(key=lambda x: re.search(date_pattern, x['url']).group())
# 计算最多能下载的文件url
max_storage = download_size  
current_size = 0
cuurent_num = 0
for entry in entries:
    current_size += entry['meta']['content_length']
    cuurent_num += entry['meta']['record_count']
    if current_size > max_storage:
        break
    urls.append(entry['url'])
print('totol num', cuurent_num)
# exit(0)
# 将url分类
works_url = []
authors_url = []
concepts_url = []
insititutions_url = []
venues_url = []
for url in urls:
    if('works' in url):
        works_url.append(url)
    elif('authors' in url):
        authors_url.append(url)
    elif('concepts' in url):
        concepts_url.append(url)
    elif('institutions' in url):
        insititutions_url.append(url)
    elif('venues' in url):
        venues_url.append(url)
urls_path = os.path.join('.', 'urls')
# 前缀
prefix = 'https://openalex.s3.amazonaws.com/'
url_pattern = r'data*'
# 先清除url_path中之前的文件
for directory, subdir_list, file_list in os.walk(urls_path):
    for file in file_list:
        os.remove(os.path.join(urls_path, file))
# 将url写入文件
with open(os.path.join(urls_path, 'works_url.txt'), 'w') as f:
    for url in works_url:
        f.write(prefix + url[url.find('data'):] + '\n')
with open(os.path.join(urls_path, 'authors_url.txt'), 'w') as f:
    for url in authors_url:
        f.write(prefix + url[url.find('data'):] + '\n')
with open(os.path.join(urls_path, 'concepts_url.txt'), 'w') as f:
    for url in concepts_url:
        f.write(prefix + url[url.find('data'):] + '\n')
with open(os.path.join(urls_path, 'institutions_url.txt'), 'w') as f:
    for url in insititutions_url:
        f.write(prefix + url[url.find('data'):] + '\n')
with open(os.path.join(urls_path, 'venues_url.txt'), 'w') as f:
    for url in venues_url:
        f.write(prefix + url[url.find('data'):] + '\n')
