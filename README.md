# Insektensommertools

## Building

```shell
go build
```

## Using

```shell
./insektensommertools -config juni2020-06.env
```

## Example crontab entry

```cron
*/15 * * * * cd /home/intevation/insekten && ./insektensommertools -config august2020-08.env > /tmp/cron_insekten.log 2>&1
```

## Campain dates

- 01.06.2018-18.06.2018
- 01.08.2018-18.08.2018
- 31.05.2019-16.06.2019
- 02.08.2019-18.08.2019
- 29.05.2020-14.06.2020
- 31.07.2020-16.08.2020
