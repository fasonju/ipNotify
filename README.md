# ipNotify

ipNotify is an event triggering tool for when the host's public ip changes.

## Deployment

To deploy this project copy and modify the `docker-compose.yml` file and run

```bash
  docker compose up -d
```

Alternatively you can simply run the binary alongside a .env file or by setting env variables manually

## Environment Variables

Configure the behavior of the public IP tracker using the following environment variables:

- `INTERVAL`  
  Interval between IP checks. Supports duration formats like `1s`, `5m`, `1h`, `1d`, or `1w` (seconds, minutes, hours, days, weeks).  
  Example: `INTERVAL=10s`

- `IPV4_ENABLED`  
  Enable tracking of the public IPv4 address.  
  Values: `true` or `false`  
  Example: `IPV4_ENABLED=true`

- `IPV6_ENABLED`  
  Enable tracking of the public IPv6 address.  
  Values: `true` or `false`  
  Example: `IPV6_ENABLED=true`

- `SMTP_ENABLED`  
  Enable email notifications when the IP address changes.  
  Values: `true` or `false`  
  Example: `SMTP_ENABLED=true`

- `SMTP_SERVER`  
  SMTP server used for sending notification emails.  
  Example: `SMTP_SERVER=smtp.gmail.com`

- `SMTP_PORT`  
  Port used by the SMTP server.  
  Example: `SMTP_PORT=587`

- `SMTP_USERNAME`  
  Username for authenticating with the SMTP server (usually your email address).  
  Example: `SMTP_USERNAME=example@gmail.com`

- `SMTP_PASSWORD`  
  Password or app-specific password for SMTP authentication.  
  Example: `SMTP_PASSWORD=your_app_password`

- `SMTP_FROM`  
  The email address the notification will be sent from.  
  Example: `SMTP_FROM=example@gmail.com`

- `SMTP_TO`  
  The recipient email address for IP change notifications.  
  Example: `SMTP_TO=example@gmail.com`

- `SCRIPTS_ENABLED`
  Enable tracking of the public IPv6 address.  
  Values: `true` or `false`  
  Example: `IPV6_ENABLED=true`

## Contributing

Contributions are always welcome!

## License

[MIT](https://choosealicense.com/licenses/mit/)

## Features

- ipv4 and/or ipv6 tracking
- SMTP notifications
- Executing arbitrary scripts
