create or alter view dbo.Balance as
with t as (
    select max(s.SurebetId) last,
           s.ServiceName    service
    from dbo.Side s
    where s.BetStatus = 'Ok'
    group by s.ServiceName
)
select s.ServiceName, s.Balance, s.CreatedAt
from dbo.Side s
         join t on t.last = s.SurebetId and t.service = s.ServiceName
