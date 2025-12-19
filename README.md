# monoctl

**monoctl** is a Go-based CLI tool to manage **configuration**, **data sources**, and **Google Sheets integrations** from the command line.

It lets you define reusable data sources (MongoDB, SQL, etc.) and push their data directly into Google Sheets.

---

## Notes & Tips

* **Config files** are stored in the OS config directory: `~/.config/monoctl/config.yml`
* **Data sources** are saved as YAML files in: `~/.config/monoctl/data-sources/<type>/<name>.yml`
* **Sheets** require a Google service account JSON file.
  * `--rows` accepts **2D JSON arrays** for inserting multiple rows.
  * `--data-source` pulls data directly from a saved data source (MongoDB, MariaDB, etc.).

---

## Installation

### Install via binary

```bash
sudo wget -O /usr/local/bin/monoctl https://github.com/salmanwaheed/monoctl/releases/download/v0.1/monoctl-linux-x86_64
sudo chmod +x /usr/local/bin/monoctl

# verify installation
monoctl version
```

---

### Build from source

Requirements: **Go**, **Git**, **Make**.

```bash
git clone https://github.com/salmanwaheed/monoctl.git
cd monoctl
make build

# verify installation
./bin/monoctl version
```

---

## Commands

### version

Show monoctl version information.

```bash
monoctl version
monoctl version --format "{{json .}}" # JSON output
monoctl version --format "{{yaml .}}" # YAML output
monoctl version --format "{{.}}"      # default Go template
```

---

### config

Manage monoctl global configuration.

```bash
monoctl config set <key> <value>             # Set a config value
monoctl config unset <key>                   # Remove a config value
monoctl config get <key>                     # Get a config value

monoctl config list                          # List all config values
monoctl config list --format "{{json .}}"
monoctl config list --format "{{yaml .}}"
monoctl config list --format "{{.}}"
```

---

### data-source

Manage reusable data source definitions.

> Allowed operators: `eq, ne, gt, gte, lt, lte, in, nin, regex, exists`.

```bash
# Create a data source
monoctl data-source create <type>/<name> \
  --uri "mongodb://<host>:27017/<db_name>?ssl=true" \
  --table <collection> \
  --limit 10 \
  --select 'f1,f2,...' \
  --sort '{"k1":"desc|asc"}' \
  --filter '{"k1":"v1","k1":{"gte":"2025-12-17T11:00:00+04:00"}}' \
  --overwrite

monoctl data-source delete <type>/<name>            # Delete a data source
monoctl data-source view <type>/<name>              # View data source YAML
monoctl data-source view <type>/<name> --dry-run    # Execute and preview data
monoctl data-source list                            # List all data sources
```

---

### sheet

Insert rows into a Google Sheet.

#### Using a data source

```bash
monoctl sheet \
  --id <sheet_id> \
  --tab-name <sheet_name> \
  --auth /path/to/google-service-account.json \
  --data-source <type>/<name>
```

#### Using rows directly

```bash
monoctl sheet \
  --id <sheet_id> \
  --tab-name <sheet_name> \
  --auth /path/to/google-service-account.json \
  --rows '[["name","email"],["Salman","salman@example.com"]]'
```

---

## Example Workflow

### 1. Create a MongoDB data source

```bash
monoctl data-source create mongodb/users \
  --uri mongodb://localhost:27017/db \
  --table users \
  --limit 10 \
  --select 'name,email' \
  --sort '{"name":"desc"}' \
  --filter '{"active":true}'
```

---

### 2. Push data into Google Sheets

```bash
monoctl sheet \
  --id 1LEh1... \
  --tab-name Sheet1 \
  --auth google-service-account.json \
  --data-source mongodb/users
```

---

## Google Service Account Setup (Required for Sheets)

monoctl uses a **Google Service Account** to authenticate and write data to Google Sheets.

Follow these steps **once** to generate the required JSON file.

---

### 1. Create a Google Cloud Project

1. Go to **Google Cloud Console** [https://console.cloud.google.com/](https://console.cloud.google.com/)
2. Click **Select Project → New Project**
3. Enter a project name and click **Create**

---

### 2. Enable Google Sheets API

1. In the Cloud Console, open **APIs & Services → Library**
2. Search for **Google Sheets API**
3. Click **Enable**

---

### 3. Create a Service Account

1. Go to **APIs & Services → Credentials**
2. Click **Create Credentials → Service Account**
3. Enter:
  * **Name**: `monoctl-readonly-access`
  * **Role**: *(Skip or choose Viewer - not required here)*
4. Click **Done**

---

### 4. Generate Service Account Key (JSON)

1. In **Credentials**, click the service account you just created
2. Go to the **Keys** tab
3. Click **Add Key → Create new key**
4. Select **JSON**
5. Download the file and save it securely

Example:

```text
~/.config/monoctl/google-service-account.json
```

⚠️ **Keep this file secret** → it grants access to your Sheets.

---

### 5. Share Google Sheet with Service Account

1. Open your Google Sheet in a browser
2. Click **Share**
3. Copy the client email (looks like: `monoctl-readonly-access@<project_id>.iam.gserviceaccount.com`)
4. Paste it into the Share dialog
5. Set permission to **Editor**
6. Click **Send**

---

### 6. Use the JSON file with `monoctl sheet ... --auth ~/.config/monoctl/google-service-account.json`

---

## Notes & Security Tips

* You **do NOT** need OAuth or browser login
* Service accounts work **server-to-server**
* Do **NOT** commit the JSON file to Git
* Recommended permissions: **Editor only on required sheets**
* You can reuse the same service account json file for multiple sheets

---

## Troubleshooting

* **403 Permission Denied** → Sheet not shared with service account email.
* **Invalid credentials** → Wrong JSON file or corrupted key.
* **API not enabled** → Enable Google Sheets API in Cloud Console.

---

## TODO

* Refactor commands, flags, and internal structure.
* Load default command flag values from config.
* Add Google Sheets integration with MariaDB as a data source.
* Add Google Drive support for downloading and storing HTML pages.
* Add `--query-from /path/to/query.json`.
* Add `--rows-from /path/to/rows.json`.
* Support multiple `--row` flags:
  ```bash
  --row "name,email,telephone" \
  --row "salman,salman@example.com,123456789"
  ```
* Enforce `monoctl help <command>` and remove `monoctl --help`.
