# 简介
Windows小工具库，基于fyne v2实现。

## 功能列表
1. 分割文本文件(目前只支持csv,txt两种格式)
```
将大文件分割为小分文件，即linux split命令的windows版本
````
2. 计算文件Hash值(MD5,SHA1,SH256,SHA512,CRC等)

## 使用
1. 首页

![首页](/images/init.png "首页")

2. 文件分割

![文件分割](/images/file-spliter.png "分割")

3. 点击"打开"文件

![打开文件](/images/open-file.png "打开文件")

4. 选中文件

![选中文件](/images/selected-file.png "选中文件")

# 打包命令
```
fyne package -os windows -icon ./images/crocodile.png
```