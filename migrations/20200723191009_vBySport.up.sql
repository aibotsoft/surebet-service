create or alter view dbo.vBySport as
select top 1000 sport,
                count(SurebetId)                                               count,
                cast(sum(a_stake + b_stake) as int)                            depo,
                cast(sum(win) as int)                                          gross,
                cast(sum(win) * 100 / sum(a_stake + b_stake) as decimal(5, 2)) perc,
                cast(avg(profit) as decimal(5, 3))                                 avg_perc,
                avg(roi)                                                       avg_roi
from Summary
group by sport
order by 3 desc