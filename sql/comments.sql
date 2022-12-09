select
    l.c,
    t2.name fund,
    t.name,
    t.twitter_id,
    t.comment,
    t.tweet_at
from
    tweets t
    inner join tickers t2
        on t.ticker = t2.ticker
    inner join (select
                    ticker,
                    count(ticker) c
                from
                    tweets
                group by
                    ticker) l
        on t.ticker = l.ticker
order by
    l.c desc,
    t.ticker asc,
    t.tweet_at asc
