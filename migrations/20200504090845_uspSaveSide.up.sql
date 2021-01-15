create or alter proc dbo.uspSaveSide @SurebetId bigint,
                                     @SideIndex tinyint,
                                     @ServiceName varchar(1000) = null,
                                     @SportName varchar(1000) = null,
                                     @LeagueName varchar(1000) = null,
                                     @Home varchar(1000) = null,
                                     @Away varchar(1000) = null,
                                     @MarketName varchar(1000) = null,
                                     @Price decimal(9, 5) = null,
                                     @Initiator bit = null,
                                     @Starts datetimeoffset = null,
                                     @EventId varchar(1000) = null,
                                     @CheckId bigint = null,
                                     @AccountId tinyint = null,
                                     @AccountLogin varchar(1000) = null,
                                     @CheckStatus varchar(1000) = null,
                                     @StatusInfo varchar(1000) = null,
                                     @CountLine int = null,
                                     @CountEvent int = null,
                                     @AmountEvent int = null,
                                     @MinBet decimal(9, 5) = null,
                                     @MaxBet decimal(9, 3) = null,
                                     @Balance int = null,
                                     @CheckPrice decimal(9, 5) = null,
                                     @Currency decimal(9, 5) = null,
                                     @CheckDone bigint = null,
                                     @CalcStatus varchar(1000) = null,
                                     @MaxStake decimal(9, 3) = null,
                                     @MinStake decimal(9, 5) = null,
                                     @MaxWin decimal(9, 5) = null,
                                     @Stake decimal(9, 5) = null,
                                     @Win decimal(9, 5) = null,
                                     @IsFirst bit = null,
                                     @ToBetId bigint = null,
                                     @TryCount int = null,
                                     @BetStatus varchar(1000) = null,
                                     @BetStatusInfo varchar(1000) = null,
                                     @Start bigint = null,
                                     @Done bigint = null,
                                     @BetPrice decimal(9, 5) = null,
                                     @BetStake decimal(9, 5) = null,
                                     @ApiBetId varchar(1000) = null,
                                     @SubService varchar(1000) = null
as
begin
    set nocount on

    MERGE dbo.Side AS t
    USING (select @SurebetId, @SideIndex) s (SurebetId, SideIndex)
    ON (t.SurebetId = s.SurebetId and t.SideIndex = s.SideIndex)

    WHEN MATCHED THEN
        UPDATE
        SET ServiceName = @ServiceName,
            SportName = @SportName,
            LeagueName = @LeagueName,
            Home = @Home,
            Away = @Away,
            MarketName = @MarketName,
            Price = @Price,
            Initiator = @Initiator,
            Starts = @Starts,
            EventId = @EventId,
            CheckId = @CheckId,
            AccountId = @AccountId,
            AccountLogin = @AccountLogin,
            CheckStatus = @CheckStatus,
            StatusInfo = @StatusInfo,
            CountLine = @CountLine,
            CountEvent = @CountEvent,
            AmountEvent = @AmountEvent,
            MinBet = @MinBet,
            MaxBet = @MaxBet,
            Balance = @Balance,
            CheckPrice = @CheckPrice,
            Currency = @Currency,
            CheckDone = @CheckDone,
            CalcStatus = @CalcStatus,
            MaxStake = @MaxStake,
            MinStake = @MinStake,
            MaxWin = @MaxWin,
            Stake = @Stake,
            Win = @Win,
            IsFirst = @IsFirst,
            ToBetId = @ToBetId,
            TryCount = @TryCount,
            BetStatus = @BetStatus,
            BetStatusInfo = @BetStatusInfo,
            Start = @Start,
            Done = @Done,
            BetPrice = @BetPrice,
            BetStake = @BetStake,
            ApiBetId = @ApiBetId,
            SubService = @SubService,
            UpdatedAt   = sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (SurebetId, SideIndex, ServiceName, SportName, LeagueName, Home, Away, MarketName, Price,
                Initiator, Starts, EventId, CheckId, AccountId, AccountLogin, CheckStatus, StatusInfo,
                CountLine, CountEvent, AmountEvent, MinBet, MaxBet, Balance, CheckPrice, Currency,
                CheckDone, CalcStatus, MaxStake, MinStake, MaxWin, Stake, Win, IsFirst, ToBetId, TryCount,
                BetStatus, BetStatusInfo, Start, Done, BetPrice, BetStake, ApiBetId, SubService)
        VALUES (s.SurebetId, s.SideIndex, @ServiceName, @SportName, @LeagueName, @Home, @Away, @MarketName, @Price,
                @Initiator, @Starts, @EventId, @CheckId, @AccountId, @AccountLogin, @CheckStatus, @StatusInfo,
                @CountLine, @CountEvent, @AmountEvent, @MinBet, @MaxBet, @Balance, @CheckPrice, @Currency,
                @CheckDone, @CalcStatus, @MaxStake, @MinStake, @MaxWin, @Stake, @Win, @IsFirst, @ToBetId, @TryCount,
                @BetStatus, @BetStatusInfo, @Start, @Done, @BetPrice, @BetStake, @ApiBetId, @SubService);
end


--     declare @Id bigint
--
--     select @Id = SurebetId from dbo.Side where SurebetId = @SurebetId and SideIndex = @SideIndex
--     if @@rowcount = 0
--         insert into dbo.Side (SurebetId, SideIndex, ServiceName, SportName, LeagueName, Home, Away, MarketName, Price,
--                               Initiator, Starts, EventId, CheckId, AccountId, AccountLogin, CheckStatus, StatusInfo,
--                               CountLine, CountEvent, AmountEvent, MinBet, MaxBet, Balance, CheckPrice, Currency,
--                               CheckDone, CalcStatus, MaxStake, MinStake, MaxWin, Stake, Win, IsFirst, ToBetId, TryCount,
--                               BetStatus, BetStatusInfo, Start, Done, BetPrice, BetStake, ApiBetId)
--         output inserted.SurebetId
--
--         values (@SurebetId, @SideIndex, @ServiceName, @SportName, @LeagueName, @Home, @Away, @MarketName, @Price,
--                 @Initiator, @Starts, @EventId, @CheckId, @AccountId, @AccountLogin, @CheckStatus, @StatusInfo,
--                 @CountLine, @CountEvent, @AmountEvent, @MinBet, @MaxBet, @Balance, @CheckPrice, @Currency,
--                 @CheckDone, @CalcStatus, @MaxStake, @MinStake, @MaxWin, @Stake, @Win, @IsFirst, @ToBetId, @TryCount,
--                 @BetStatus, @BetStatusInfo, @Start, @Done, @BetPrice, @BetStake, @ApiBetId )
--     else
--         select @Id