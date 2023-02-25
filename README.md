
# TGRHB - Telegram Reddit Hot Bot

A Telegram bot that sends random top posts to users on button click

## Installation


```bash
  # PREREQUESITES: docker installed and running
  #
  # 1. Clone this repo on your local machine
  # 2. Copy the app.env to app.dev.env. Add values to these variables:
  #       TGRHB_REDDIT_AUTH - see reddit api reference on how to obtain install
  #       TGRHB_TG_TOKEN - your telegram bot TGRHB_TG_TOKEN
  #       TGRHB_DB_DRIVER - set to postgres
  #       TGRHB_DB_SOURCE - pg connection string. Note that this should be identical to docker compose
  #       TGRHB_ENCRYPTION_KEY - your key to encrypt reddit token for storing it. Should be 8, 16, 24 or 32 characters (bytes) long
  # 3. Then run in project root: 
  
  docker-compose up
```
    