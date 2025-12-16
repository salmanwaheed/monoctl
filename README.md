# monoctl

just release first version for cron job.

```sh
# manual rows
monoctl sheet --id xxx --tab-name xxx --auth service-account.json --rows '[["name","email","telephone"],["salman","salman@example.com","123456789"]]'

# dynamic rows: get rows from database
monoctl sheet --id xxx --tab-name xxx --auth ~/.ssh/google-service-account.json --data-source mongodb

# create data source
monoctl data-source create mongodb/fetch-leads --uri xxx --table xxx --limit 2 --fields 'name,email,telephone' --query '{}'
```

## TODO
- Refactor code, flags, commands and etc.
- Pass `--rows-file /path/data.json`.
- Pass multiple `--row "name,email,telephone" --row "salman,salman@example.com,123456789"` flags.
- Use `monoctl help command`, and remove `monoctl --help`.
- `monoctl data-source view <type>/<name> --dry-run`.
