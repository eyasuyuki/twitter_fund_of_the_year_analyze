WITH CandidateVote AS (
    select count(tweets.ticker) c,
           tickers.ticker ticker
    from tweets
             inner join tickers
                        on tweets.ticker = tickers.ticker
    group by tweets.ticker
),
     RankedVotes AS (
         SELECT
             ticker,
             c,
             RANK() OVER (ORDER BY c DESC) AS r
         FROM
             CandidateVote
     )
SELECT
    r,
    c,
    ticker
FROM
    RankedVotes;

