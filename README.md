# monoctl

monoctl is a Go based CLI tool to manage configuration, data sources, and Google Sheets integrations.

---

## Installation

Build from source, if you have go installed:
```bash
git clone https://github.com/salmanwaheed/monoctl.git

cd monoctl
make build

# verify installation
./bin/monoctl version
```

Install via Binary file:
```bash
sudo wget -O /usr/local/bin/monoctl https://github.com/salmanwaheed/monoctl/releases/download/v0.1/monoctl-linux-x86_64
sudo chmod +x /usr/local/bin/monoctl

# verify installation
monoctl version
```

---

## Example Workflow

1. Create a MongoDB data source:

```bash
monoctl data-source create mongodb/users --uri mongodb://localhost:27017/db --table users --limit 10 --fields 'name,email' --query '{"active":true}'
```

2. List all data sources:

```bash
monoctl data-source list
```

3. Push data from data source to a Google Sheet:

```bash
monoctl sheet --id 1LEh1... --tab-name Sheet1 --auth google-service-account.json --data-source mongodb/users
```

4. Append custom rows:

```bash
monoctl sheet --id 1LEh1... --tab-name Sheet1 --auth google-service-account.json --rows '[["name","email"],["Salman","salman@example.com"]]'
```

---

## TODO
- Refactor code, flags, commands and etc.
- Pass `--rows-file /path/data.json`.
- Pass multiple `--row "name,email,telephone" --row "salman,salman@example.com,123456789"` flags.
- Use `monoctl help command`, and remove `monoctl --help`.
