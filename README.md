# getfile
拉取文件到本地，在打包机a上部署此程序，测试机b拉取打包机a中指定文件夹下的文件

编译 + 压缩：

linux
go build -ldflags="-s -w" -o getfile main.go && upx -9 getfile

window
go build -ldflags="-s -w" -o getfile.exe main.go && upx -9 getfile.exe


# 安装upx
cd /usr/local

wget -c https://github.com/upx/upx/releases/download/v3.96/upx-3.96-amd64_linux.tar.xz && tar xvf upx-3.96-amd64_linux.tar.xz

# 删除、重命名
rm -rf upx-3.96-amd64_linux.tar.xz && mv upx-3.96-amd64_linux upx

cd ~

vi .bashrc

alias upx='/usr/local/upx/upx'

source .bashrc


