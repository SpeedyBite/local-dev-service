TODO:
- Dump database if want to use database local
- Replace CHANNEL_DB_* by new database configuration
- Load Vault for environment
- EXPORT ALL to .<environemt>.env file. Except following attr:

  + CHANNEL_DB_* (if would like to use local)
  + <dependency>_URL (replace by the values in the configuration)

- Generate docker-compose file for service

2. Write function to execute a jobs