set -e  # 出错后退出 shell
set -x  # 打开调试 shell 命令

remote="texas-auth@118.24.147.175:~/cpay"

env GOOS=linux GOARCH=amd64 go build -o bin/cpay
rsync bin/cpay $remote

rsync -r config/default.yaml $remote/config/default.yaml
rsync -r config/gamma.yaml $remote/config/gamma.yaml
