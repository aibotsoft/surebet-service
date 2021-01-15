create or alter view dbo.vByServices as
select top 1000 dbo.fnSortedConcat(a_service, b_service) services,
                count(SurebetId)                                               count,
                cast(sum(a_stake + b_stake) as int)                            depo,
                cast(sum(win) as int)                                          gross,
                cast(sum(win) * 100 / sum(a_stake + b_stake) as decimal(5, 2)) perc,
                cast(avg(profit) as decimal(5, 2))                             avg_perc,
                avg(roi)                                                       avg_roi
from Summary
group by dbo.fnSortedConcat(a_service, b_service)
order by 3 desc

