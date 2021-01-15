create or alter view dbo.vPinSub as
select top 100 cast(sum(daf_loss) as int)                                    daf_loss,
               cast(sum(other_loss) as int)                                  other_loss,
               cast(sum(stake_sum) as int)                                   stake_sum,
               cast(sum(diff) as int)                                        diff,
               cast(sum(daf_loss) * 100.0 / sum(stake_sum) as decimal(5, 2)) skew_percent,
               sum(count)                                                    count,
               b_sub
from (
         select sum(a_win_loss)                   daf_loss,
                sum(b_win_loss)                   other_loss,
                sum(a_win_loss) + sum(b_win_loss) diff,
                sum(a_stake + b_stake)            stake_sum,
                count(a_win_loss)                 count,
                b_sub
         from Summary
         where a_service = 'Pinnacle'
           and win is not null
         group by b_sub
         union
         select sum(b_win_loss)                   daf_loss,
                sum(a_win_loss)                   other_loss,
                sum(b_win_loss) + sum(a_win_loss) diff,
                sum(a_stake + b_stake)            stake_sum,
                count(b_win_loss)                 count,
                a_sub
         from Summary
         where b_service = 'Pinnacle'
           and win is not null
         group by a_sub) as t
group by b_sub
order by stake_sum desc;

-- create or alter view dbo.vDafSport as
select top 100 cast(sum(daf_loss) as int)                                    daf_loss,
               cast(sum(other_loss) as int)                                  other_loss,
               cast(sum(stake_sum) as int)                                   stake_sum,
               cast(sum(diff) as int)                                        diff,
               cast(sum(daf_loss) * 100.0 / sum(stake_sum) as decimal(5, 2)) skew_percent,
               sum(count)                                                    count,
               sport
from (
         select sum(a_win_loss)                   daf_loss,
                sum(b_win_loss)                   other_loss,
                sum(a_win_loss) + sum(b_win_loss) diff,
                sum(a_stake + b_stake)            stake_sum,
                count(a_win_loss)                 count,
                sport
         from Summary
         where a_service = 'Pinnacle'
           and win is not null
         group by sport
         union
         select sum(b_win_loss)                   daf_loss,
                sum(a_win_loss)                   other_loss,
                sum(b_win_loss) + sum(a_win_loss) diff,
                sum(a_stake + b_stake)            stake_sum,
                count(b_win_loss)                 count,
                sport
         from Summary
         where b_service = 'Pinnacle'
           and win is not null
         group by sport) as t
where count > 10
group by sport
order by stake_sum desc;


create or alter view dbo.vPinVsFast as
select top 100 cast(sum(daf_loss) as int)                                    pin_pl,
               cast(sum(other_loss) as int)                                  other_loss,
               cast(sum(stake_sum) as int)                                   stake_sum,
               cast(sum(diff) as int)                                        diff,
               cast(sum(daf_loss) * 100.0 / sum(stake_sum) as decimal(5, 2)) skew_percent,
               sum(count)                                                    count
from (select sum(a_win_loss)                   daf_loss,
             sum(b_win_loss)                   other_loss,
             sum(a_win_loss) + sum(b_win_loss) diff,
             sum(a_stake + b_stake)            stake_sum,
             count(a_win_loss)                 count
      from Summary
      where a_service = 'Pinnacle'
        and b_sub in ('bf', 'bdaq', 'mbook', 'isn', 'sing2')
--         and b_service = 'Pinnacle'
        and win is not null
      union
      select sum(b_win_loss)                   daf_loss,
             sum(a_win_loss)                   other_loss,
             sum(b_win_loss) + sum(a_win_loss) diff,
             sum(a_stake + b_stake)            stake_sum,
             count(b_win_loss)                 count
      from Summary
      where b_service = 'Pinnacle'
        and a_sub in ('bf', 'bdaq', 'mbook', 'isn', 'sing2')
--         and a_service = 'Pinnacle'
        and win is not null
     ) as t
order by 1 desc