sudo sed -i '' "/github/d" /etc/hosts
go run main.go
cat dns_github.txt >> /etc/hosts
dscacheutil -flushcache
