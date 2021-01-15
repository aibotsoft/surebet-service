create or alter view dbo.vDafabet as
select top 100 cast(sum(daf_loss) as int)                                    daf_loss,
               cast(sum(other_loss) as int)                                  other_loss,
               cast(sum(stake_sum) as int)                                   stake_sum,
               cast(sum(diff) as int)                                        diff,
               cast(sum(daf_loss) * 100.0 / sum(stake_sum) as decimal(5, 2)) skew_percent,
               sum(count)                                                    count,
--                cast(sum(reduced_profit) as int)                              reduced_profit,
               b_service
from (
         select sum(a_win_loss)                   daf_loss,
                sum(b_win_loss)                   other_loss,
                sum(a_win_loss) + sum(b_win_loss) diff,
                sum(a_stake + b_stake)            stake_sum,
                count(a_win_loss)                 count,
--                 sum(a_win_loss / a_stake)         reduced_profit,

                b_service
         from Summary
         where a_service = 'Dafabet'
           and win is not null
         group by b_service
         union
         select sum(b_win_loss)                   daf_loss,
                sum(a_win_loss)                   other_loss,
                sum(b_win_loss) + sum(a_win_loss) diff,
                sum(a_stake + b_stake)            stake_sum,
                count(b_win_loss)                 count,
--                 sum(b_win_loss / b_stake)         reduced_profit,

                a_service
         from Summary
         where b_service = 'Dafabet'
           and win is not null
         group by a_service) as t
group by b_service
order by stake_sum desc;

create or alter view dbo.vDafSub as
select top 100 cast(sum(daf_loss) as int)                                    daf_loss,
               cast(sum(other_loss) as int)                                  other_loss,
               cast(sum(stake_sum) as int)                                   stake_sum,
               cast(sum(diff) as int)                                        diff,
               cast(sum(daf_loss) * 100.0 / sum(stake_sum) as decimal(5, 2)) skew_percent,
               sum(count)                                                    count,
--                cast(sum(reduced_profit) as int)                              reduced_profit,
               b_sub
from (
         select sum(a_win_loss)                   daf_loss,
                sum(b_win_loss)                   other_loss,
                sum(a_win_loss) + sum(b_win_loss) diff,
                sum(a_stake + b_stake)            stake_sum,
                count(a_win_loss)                 count,
--                 sum(a_win_loss / a_stake)         reduced_profit,
                b_sub
         from Summary
         where a_service = 'Dafabet'
           and win is not null
         group by b_sub
         union
         select sum(b_win_loss)                   daf_loss,
                sum(a_win_loss)                   other_loss,
                sum(b_win_loss) + sum(a_win_loss) diff,
                sum(a_stake + b_stake)            stake_sum,
                count(b_win_loss)                 count,
--                 sum(b_win_loss / b_stake)         reduced_profit,

                a_sub
         from Summary
         where b_service = 'Dafabet'
           and win is not null
         group by a_sub) as t
group by b_sub
order by stake_sum desc;

create or alter view dbo.vDafWeek as
select top 100 week,
               cast(sum(daf_loss) as int)                                    daf_loss,
               cast(sum(other_loss) as int)                                  other_loss,
               cast(sum(stake_sum) as int)                                   stake_sum,
               cast(sum(diff) as int)                                        diff,
               cast(sum(daf_loss) * 100.0 / iif(sum(stake_sum)=0, 1, sum(stake_sum)) as decimal(5, 2)) skew_percent,
               sum(count)                                                    count,
               avg(avg_profit)                                               avg_profit
--                cast(sum(reduced_profit) as int)                              reduced_profit

from (select DATEPART(Week, CreatedAt)         week,
             sum(a_win_loss)                   daf_loss,
             sum(b_win_loss)                   other_loss,
             sum(a_win_loss) + sum(b_win_loss) diff,
             sum(a_stake + b_stake)            stake_sum,
             count(a_win_loss)                 count,
             avg(real_profit)                  avg_profit
--              sum(a_win_loss / a_stake)         reduced_profit
--           var FastServiceList = []string{"bf", "mbook", "bdaq", "pin", "pin88", "isn", "sing2", "penta88"}

      from Summary
      where a_service = 'Dafabet'
        and b_sub in ('pin', 'pin88', 'bf', 'bdaq', 'mbook', 'isn', 'sing2', 'penta88')
--         and b_service = 'Pinnacle'
        and win is not null
      group by DATEPART(Week, CreatedAt)
      union
      select DATEPART(Week, CreatedAt)         week,
             sum(b_win_loss)                   daf_loss,
             sum(a_win_loss)                   other_loss,
             sum(b_win_loss) + sum(a_win_loss) diff,
             sum(a_stake + b_stake)            stake_sum,
             count(b_win_loss)                 count,
             avg(real_profit)                  avg_profit
--              sum(b_win_loss / b_stake)         reduced_profit
      from Summary
      where b_service = 'Dafabet'
        and a_sub in ('pin', 'pin88', 'bf', 'bdaq', 'mbook', 'isn', 'sing2', 'penta88')
--         and a_service = 'Pinnacle'
        and win is not null
      group by DATEPART(Week, CreatedAt)
     ) as t
group by week
order by 1 desc;

create or alter view dbo.vDafDay as
select top 1000 day,
               cast(sum(daf_loss) as int)                                    daf_loss,
               cast(sum(other_loss) as int)                                  other_loss,
               cast(sum(stake_sum) as int)                                   stake_sum,
               cast(sum(diff) as int)                                        diff,
               cast(sum(daf_loss) * 100.0 / iif(sum(stake_sum)=0, 1, sum(stake_sum)) as decimal(5, 2)) skew_percent,
               sum(count)                                                    count,
               avg(avg_profit)                                               avg_profit
from (select cast(CreatedAt as date)        day,
             sum(a_win_loss)                   daf_loss,
             sum(b_win_loss)                   other_loss,
             sum(a_win_loss) + sum(b_win_loss) diff,
             sum(a_stake + b_stake)            stake_sum,
             count(a_win_loss)                 count,
             avg(real_profit)                  avg_profit
      from Summary
      where a_service = 'Dafabet'
        and b_sub in ('pin', 'pin88', 'bf', 'bdaq', 'mbook', 'isn', 'sing2', 'penta88')
        and win is not null
      group by cast(CreatedAt as date)
      union
      select cast(CreatedAt as date)         day,
             sum(b_win_loss)                   daf_loss,
             sum(a_win_loss)                   other_loss,
             sum(b_win_loss) + sum(a_win_loss) diff,
             sum(a_stake + b_stake)            stake_sum,
             count(b_win_loss)                 count,
             avg(real_profit)                  avg_profit
      from Summary
      where b_service = 'Dafabet'
        and a_sub in ('pin', 'pin88', 'bf', 'bdaq', 'mbook', 'isn', 'sing2', 'penta88')
        and win is not null
      group by cast(CreatedAt as date)
     ) as t
group by day
order by 1 desc;

create or alter view dbo.vDafWD as
select top 1000 weekday,
               cast(sum(daf_loss) as int)                                    daf_loss,
               cast(sum(other_loss) as int)                                  other_loss,
               cast(sum(stake_sum) as int)                                   stake_sum,
               cast(sum(diff) as int)                                        diff,
               cast(sum(daf_loss) * 100.0 / iif(sum(stake_sum)=0, 1, sum(stake_sum)) as decimal(5, 2)) skew_percent,
               sum(count)                                                    count,
               avg(avg_profit)                                               avg_profit
from (select dbo.fnRuWeekDay(DATEPART(weekday , CreatedAt))     weekday,
             sum(a_win_loss)                   daf_loss,
             sum(b_win_loss)                   other_loss,
             sum(a_win_loss) + sum(b_win_loss) diff,
             sum(a_stake + b_stake)            stake_sum,
             count(a_win_loss)                 count,
             avg(real_profit)                  avg_profit
      from Summary
      where a_service = 'Dafabet'
        and b_sub in ('pin', 'pin88', 'bf', 'bdaq', 'mbook', 'isn', 'sing2', 'penta88')
        and win is not null
      group by dbo.fnRuWeekDay(DATEPART(weekday , CreatedAt))
      union
      select dbo.fnRuWeekDay(DATEPART(weekday , CreatedAt))       weekday,
             sum(b_win_loss)                   daf_loss,
             sum(a_win_loss)                   other_loss,
             sum(b_win_loss) + sum(a_win_loss) diff,
             sum(a_stake + b_stake)            stake_sum,
             count(b_win_loss)                 count,
             avg(real_profit)                  avg_profit
      from Summary
      where b_service = 'Dafabet'
        and a_sub in ('pin', 'pin88', 'bf', 'bdaq', 'mbook', 'isn', 'sing2', 'penta88')
        and win is not null
      group by dbo.fnRuWeekDay(DATEPART(weekday , CreatedAt))
     ) as t
group by weekday
order by 2;


----------------------------------------------------------------------------------
create or alter view dbo.vDafSport as
select top 100 cast(sum(daf_loss) as int)                                    daf_loss,
               cast(sum(other_loss) as int)                                  other_loss,
               cast(sum(stake_sum) as int)                                   stake_sum,
               cast(sum(diff) as int)                                        diff,
               cast(sum(daf_loss) * 100.0 / sum(stake_sum) as decimal(5, 2)) skew_percent,
               sum(count)                                                    count,
--                cast(sum(reduced_profit) as int)                              reduced_profit,
               sport
from (
         select sum(a_win_loss)                   daf_loss,
                sum(b_win_loss)                   other_loss,
                sum(a_win_loss) + sum(b_win_loss) diff,
                sum(a_stake + b_stake)            stake_sum,
                count(a_win_loss)                 count,
--                 sum(a_win_loss / a_stake)         reduced_profit,
                sport
         from Summary
         where a_service = 'Dafabet'
           and win is not null
         group by sport
         union
         select sum(b_win_loss)                   daf_loss,
                sum(a_win_loss)                   other_loss,
                sum(b_win_loss) + sum(a_win_loss) diff,
                sum(a_stake + b_stake)            stake_sum,
                count(b_win_loss)                 count,
--                 sum(b_win_loss / b_stake)         reduced_profit,
                sport
         from Summary
         where b_service = 'Dafabet'
           and win is not null
         group by sport) as t
group by sport
order by stake_sum desc



create or alter function dbo.fnRuWeekDay(@WeekDay int) returns varchar(16) as
begin
    if @WeekDay = 2
        return 'понедельник'
    if @WeekDay = 3
        return 'вторник'
    if @WeekDay = 4
        return 'среда'
    if @WeekDay = 5
        return 'четверг'
    if @WeekDay = 6
        return 'пятница'
    if @WeekDay = 7
        return 'суббота'
    if @WeekDay = 1
        return 'воскресение'
    return 'неизвестно'
end

select DATEPART(weekday , sysdatetimeoffset())
select dbo.fnRuWeekDay(DATEPART(weekday , sysdatetimeoffset()))
------------------------------------------------------------------------------------
select DATEPART(weekday , CreatedAt)         week,
       sum(a_win_loss)                   daf_loss,
       sum(b_win_loss)                   other_loss,
       sum(a_win_loss) + sum(b_win_loss) diff,
       sum(a_stake + b_stake)            stake_sum,
       count(a_win_loss)                 count,
       avg(profit),
       avg(real_profit)
from Summary
where a_service = 'Dafabet'
  and b_service = 'Pinnacle'
  and win is not null
group by DATEPART(weekday, CreatedAt)
order by 1 desc;

select cast(CreatedAt as date)         week,
       sum(a_win_loss)                   daf_loss,
       sum(b_win_loss)                   other_loss,
       sum(a_win_loss) + sum(b_win_loss) diff,
       sum(a_stake + b_stake)            stake_sum,
       count(a_win_loss)                 count
--        sum(a_win_loss / a_stake)

from Summary
where a_service = 'Dafabet'
  and b_service = 'Pinnacle'
  and win is not null
group by cast(CreatedAt as date)
order by 1 desc;

select cast(sysdatetimeoffset() as date)
-- select top 100 DATEPART(Year, s.CreatedAt)                                year,
--                DATEPART(Month, s.CreatedAt)                               month,
--                DATEPART(Week, s.CreatedAt)                                week,
--                count(distinct b.SurebetId)                                count,
--                cast(sum(b.WinLoss) as int)                                gross,
--                cast(sum(b.Stake) as int)                                  depo,
--                cast(sum(b.WinLoss) * 100 / sum(b.Stake) as decimal(9, 2)) perc
-- from dbo.BetList b
--          join dbo.Side s on s.SurebetId = b.SurebetId and s.SideIndex = b.SideIndex
-- group by DATEPART(Year, s.CreatedAt), DATEPART(Month, s.CreatedAt), DATEPART(Week, s.CreatedAt)
-- order by 3 desc;

