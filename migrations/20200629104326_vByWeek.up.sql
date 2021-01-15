create or alter view dbo.vByWeek as
select top 100
               DATEPART(Week, s.CreatedAt)                               week,
               count(distinct b.SurebetId)                                count,
               cast(sum(b.WinLoss) as int)                                gross,
               cast(sum(b.Stake) as int)                                  depo,
               cast(sum(b.WinLoss) * 100 / sum(b.Stake) as decimal(9, 2)) perc
from dbo.BetList b
         join dbo.Side s on s.SurebetId = b.SurebetId and s.SideIndex = b.SideIndex
group by DATEPART(Week, s.CreatedAt)
order by 1 desc;