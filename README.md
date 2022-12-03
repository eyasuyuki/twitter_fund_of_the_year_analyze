twitter_fund_of_the_year_analyze
====

# execute

```shell
go run main.go
```

# SQLite

```shell
sqlite3 foy2022.db
```

## csv mode

```sqlite
.headers on
.mode csv
.once './lanking.csv'
```

## count by ticker

```sql
select
    c,
    name
from
    (select
        count(ticker) c,
        ticker
     from
        tweets
     group by
        ticker) s1,
    tickers
where
    s1.ticker=tickers.ticker
order by
    c desc
```
