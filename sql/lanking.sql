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
