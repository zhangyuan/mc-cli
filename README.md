# mc-helper

## Usage

Create the `.env` file with the sample content as below:

```env
DSN=http://ACCESS_KEY_ID:ACCESS_KEY_SECRET@service.cn-beijing.maxcompute.aliyun.com/api?project=project_name
```

If multiple SQL statements will be executed, `odps.sql.submit.mode=script` must be set. So the DSN would be like:

```env
DSN=http://ACCESS_KEY_ID:ACCESS_KEY_SECRET@service.cn-beijing.maxcompute.aliyun.com/api?project=project_name&odps.sql.submit.mode=script
```
