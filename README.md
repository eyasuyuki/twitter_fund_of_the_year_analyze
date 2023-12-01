twitter_fund_of_the_year_analyze
====

# execute

```shell
go run main.go
```

# SQLite

```shell
sqlite3 foy2023.db
```

## csv mode

```sqlite
.headers on
.mode csv
.once './ranking.csv'
```

## count by ticker

```sql
select
    count(tweets.ticker) c,
    tickers.name
from
    tweets
    inner join tickers
        on tweets.ticker=tickers.ticker
group by
    tweets.ticker
order by
    c desc;
```

File ```ranking.csv``` will be saved.
