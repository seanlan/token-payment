CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go && \
scp main fish:/data/tokenpay/main.new && \
ssh fish "supervisorctl stop tokenpay-cron && \
          \cp /data/tokenpay/main /data/tokenpay/main.last && \
          \cp /data/tokenpay/main.new /data/tokenpay/main && \
          supervisorctl start tokenpay-cron" && \
rm -rf main