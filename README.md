# Insektensommertools

- [Insektensommertools](#insektensommertools)
  - [Building](#building)
  - [Using](#using)
  - [Example env](#example-env)
  - [Example crontab entry](#example-crontab-entry)
  - [Campain dates](#campain-dates)

## Building

```shell
go build
```

## Using

```shell
./insektensommertools -config juni2020-06.env
```

## Example env

```env
AKTION="beobachtungen-2020-08"
URL="https://naturgucker.de/mobil/?modul=beobachtungenNABU&email1=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX&md5=XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX&offset=0&zeilen=0&service=XXXXXXXXXX&datumvon=31.07.2020&datumbis=16.08.2020"
FTP_HOST=XX.XXX.XXX.XXX:21
FTP_USER=XXXXXXXXX
FTP_PASSWORD=XXXXXXXXXXXXXXXXXX
FTP_PATH=insektensommer/data/
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
