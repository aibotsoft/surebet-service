create or alter view dbo.vByDay as
select top 1000 cast(s.CreatedAt AS date)                                  day,
                count(distinct b.SurebetId)                                count,
                cast(sum(b.WinLoss) as int)                                gross,
                cast(sum(b.Stake) as int)                                  depo,
                cast(sum(b.WinLoss) * 100 / sum(b.Stake) as decimal(9, 2)) perc
from dbo.BetList b
         join dbo.Side s on s.SurebetId = b.SurebetId and s.SideIndex = b.SideIndex
group by CAST(s.CreatedAt AS date)
order by 1 desc;


create or alter view dbo.vFastByDay as
select top 1000 cast(CreatedAt AS date)   day,
                count(distinct SurebetId) count,
                cast(sum(win) as int)     gross
from Summary
where a_service not in ('Dafabet', 'Sbobet')
  and b_service not in ('Dafabet', 'Sbobet')
group by CAST(CreatedAt AS date)
order by 1 desc;

create or alter view dbo.vFastByWeek as
select top 1000 datepart(Week, CreatedAt) week,
                count(distinct SurebetId) count,
                cast(sum(win) as int)     gross
from Summary
where a_service not in ('Dafabet', 'Sbobet')
  and b_service not in ('Dafabet', 'Sbobet')
group by datepart(Week, CreatedAt)
order by 1 desc;

create or alter view dbo.vFastByMonth as
select top 1000 DATEPART(Year, CreatedAt)  year,
                DATEPART(Month, CreatedAt) month,
                count(distinct SurebetId)  count,
                cast(sum(win) as int)      gross
from Summary
where a_service not in ('Dafabet', 'Sbobet')
  and b_service not in ('Dafabet', 'Sbobet')
group by datepart(Year, CreatedAt), datepart(Month, CreatedAt)
order by 1,2 desc;

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