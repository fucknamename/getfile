# getfile
拉取文件到本地，在打包机a上部署此程序，测试机b拉取打包机a中指定文件夹下的文件

★★★注意：要打包的文件权限需要配置好，不然压缩会提示失败，权限不足★★★

# 使用方式

使用wget方式  
wget http://xx.xx.xx.xx:1688/getfile/目录/这是token

使用curl方式  
curl -o 保存文件名.zip http://xx.xx.xx.xx:1688/getfile/目录/这是token

针对https的忽略处理：  
wget  ... --no-check-certificate  
curl ... -k


# 编译 + 压缩：

linux  
go build -ldflags="-s -w" -o getfile main.go && upx -9 getfile

window  
go build -ldflags="-s -w" -o getfile.exe main.go && upx -9 getfile.exe

window 下无边框  
go build -ldflags="-H windowsgui -w -s" -o getfile.exe  main.go && upx -9 getfile.exe

# windows下打包带图标
创建main.rc文件，写入IDI_ICON ICON "favicon.ico"  

在项目根目录下打开 cmd 窗口运行下面的命令生成 main.syso 文件  
windres -o main.syso main.rc  

再执行 go build -ldflags="-s -w" && upx -9 getfile.exe

# 安装upx
cd /usr/local

wget -c https://github.com/upx/upx/releases/download/v3.96/upx-3.96-amd64_linux.tar.xz && tar xvf upx-3.96-amd64_linux.tar.xz

删除、重命名

rm -rf upx-3.96-amd64_linux.tar.xz && mv upx-3.96-amd64_linux upx

cd ~

vi .bashrc

alias upx='/usr/local/upx/upx'

source .bashrc


