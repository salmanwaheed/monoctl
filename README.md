# monoctl

just release first version for cron job.

```sh
# manual rows
monoctl sheet --id xxx --tab-name xxx --auth service-account.json --rows '[["name","email","telephone"],["salman","salman@example.com","123456789"]]'

# dynamic rows: get rows from database
monoctl sheet --id xxx --tab-name xxx --auth ~/.ssh/google-service-account.json --data-source mongodb
```

## TODO
- Refactor code, flags, commands and etc.
