### 项目功能

* 解析m3u8文件，返回ts列表

#### 需求
* 上传txt文件(gbk,utf8)，文件大小约为1万行，内容为m3u8列表，每行一条记录
* 每个m3u8文件，解析为对应ts文件列表，一个m3u8可能对应几千个ts文件，所以返回文件可能是上传文件的上千倍，推荐压缩一下下载

### 
* https://github.com/viki-org/dnscache