WITH CandidateVote AS (
    select count(tweets.ticker) c,
          tickers.name name
   from tweets
            inner join tickers
                       on tweets.ticker = tickers.ticker
   group by tweets.ticker
),
RankedVotes AS (
    SELECT
        name,
        c,
        RANK() OVER (ORDER BY c DESC) AS r
    FROM
        CandidateVote
)
SELECT
    r,
    c,
    name
FROM
    RankedVotes;
