create or alter view dbo.vByMonth as
select top 100 DATEPART(Year, s.CreatedAt)                                year,
               DATEPART(Month, s.CreatedAt)                               month,
               count(distinct b.SurebetId)                                count,
               cast(sum(b.WinLoss) as int)                                gross,
               cast(sum(b.Stake) as int)                                  depo,
               cast(sum(b.WinLoss) * 100 / sum(b.Stake) as decimal(9, 2)) perc
from dbo.BetList b
         join dbo.Side s on s.SurebetId = b.SurebetId and s.SideIndex = b.SideIndex
group by DATEPART(Year, s.CreatedAt), DATEPART(Month, s.CreatedAt)
order by 2 desc;


-- select cast(sysdatetimeoffset() AS year), MONTH(sysdatetimeoffset());
-- SELECT DATEPART(Year, sysdatetimeoffset()) Year, DATEPART(Month, sysdatetimeoffset())
-- select datepart(week , sysdatetimeoffset())
