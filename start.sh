go run -race . >> output.txt &
>
code output.txt .
sleep 5 &&
# curl -X POST http://localhost:8080/back-link -d '{"comp_domains":["homedepot.com","ajmadison.com","lowes.com","bestbuy.com","build.com"],"industry":"intitle:buy home appliance inurl:forum OR inurl:forums OR inurl:discussion OR inurl:discussions OR inurl:thread OR inurl:threads"}'
curl -X POST http://localhost:8080/back-link -d '{"comp_domains":["homedepot.com","ajmadison.com","lowes.com","bestbuy.com","build.com"],"industry":"home%20appliance%20forums%20-site:doityourself.com"}'
