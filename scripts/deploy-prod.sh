set -e  # 出错后退出 shell
set -x  # 打开调试 shell 命令

remote="root@207.246.104.70:/root/cpay"

env GOOS=linux GOARCH=amd64 go build -o bin/cpay
rsync bin/cpay $remote

rsync -r config/default.yaml $remote/config/default.yaml
rsync -r config/prod.yaml $remote/config/prod.yaml
