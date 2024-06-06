# mc-cli

## Usage

### Setup the credentials

Create the `.env` file with the sample content as below:

```env
DSN=http://ACCESS_KEY_ID:ACCESS_KEY_SECRET@service.cn-beijing.maxcompute.aliyun.com/api?project=project_name
```

If multiple SQL statements will be executed, `odps.sql.submit.mode=script` must be set. So the DSN would be like:

```env
DSN=http://ACCESS_KEY_ID:ACCESS_KEY_SECRET@service.cn-beijing.maxcompute.aliyun.com/api?project=project_name&odps.sql.submit.mode=script
```

### Query the data via single SQL

```bash
./mc-cli query -s "SELECT 1"
```

or use the SQL from file:

```bash
./mc-cli query -f sample.sql
```

Variables in the same style as Dataworks ODPS node are supported via `--dataworks-vars|-v` or `dataworks-vars-file|-d`. E.g.

```bash
./mc-cli query -f sample.sql -d sample.yaml
```

### Execute SQL statement(s)

```bash
./mc-cli exec -f sample.sql
```

or 

```bash
./mc-cli exec -f sample.sql -d sample.yaml
```
