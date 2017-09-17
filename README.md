# go-godaddy-dyndns

Update DNS record with current external IP.
Domain + Record must existing (create it before executing the script)

## Build
go build

or if building for a different architecture i.e. freebsd
env GOOS=freebsd GOARCH=amd64 go build

## Usage

go-godaddy-dyndns -h

Usage of ./go-godaddy-dyndns:
  -dwr string
        GoDaddy Domain (JSON: {"domain": "yourdomain.com","record": "yourdnsrecord","type": "A"}'
  -key string
        GoDaddy API-Key
  -secret string
        GoDaddy API-Secret

Create a cronjob and execute frequently. Update of the record only happens if current external ip is different to the one stored in the domain record.

### Example
./go-godaddy-dyndns -key=<godaddy-key> -secret=<godaddy-secret> -dwr='{"domain": "<domain>","record": "<record>","type": "<record type i.e. A>"}'


