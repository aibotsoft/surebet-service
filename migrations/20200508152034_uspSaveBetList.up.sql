create or alter proc dbo.uspSaveBetList @TVP dbo.BetListType READONLY as
begin
    set nocount on

    MERGE dbo.BetList AS t
    USING @TVP s
    ON (t.SurebetId = s.SurebetId and t.SideIndex = s.SideIndex)

    WHEN MATCHED THEN
        UPDATE
        SET Price        = s.Price,
            BetId        = s.BetId,

            Stake        = s.Stake,
            WinLoss      = s.WinLoss,
            ApiBetId     = s.ApiBetId,
            ApiBetStatus = s.ApiBetStatus,
            UpdatedAt    =sysdatetimeoffset()

    WHEN NOT MATCHED THEN
        INSERT (SurebetId, SideIndex, BetId, Price, Stake, WinLoss, ApiBetId, ApiBetStatus)
        VALUES (s.SurebetId, s.SideIndex, s.BetId, s.Price, s.Stake, s.WinLoss, s.ApiBetId, s.ApiBetStatus);
end


create or alter proc dbo.uspSaveBetListNew @SurebetId bigint,
                                           @SideIndex tinyint,
                                           @BetId bigint,
                                           @Price decimal(9, 5) = null,
                                           @Stake decimal(9, 5)= null,
                                           @WinLoss decimal(9, 5)= null,
                                           @ApiBetId varchar(1000)= null,
                                           @ApiBetStatus varchar(1000)= null
as
begin
    set nocount on

    declare @Id bigint
    select @Id = SurebetId from dbo.BetList where SurebetId = @SurebetId and BetId = @BetId
    if @@rowcount = 0
        insert into dbo.BetList (SurebetId, BetId, Price, Stake, WinLoss, ApiBetId, ApiBetStatus, SideIndex)
        values (@SurebetId, @BetId, @Price, @Stake, @WinLoss, @ApiBetId, @ApiBetStatus, @SideIndex)
    else
        UPDATE dbo.BetList
        SET Price        = @Price,
--             BetId        = @BetId,
--             SideIndex    = @SideIndex,
            Stake        = @Stake,
            WinLoss      = @WinLoss,
            ApiBetId     = @ApiBetId,
            ApiBetStatus = @ApiBetStatus,
            UpdatedAt    =sysdatetimeoffset()
        where SurebetId = @SurebetId
          and BetId = @BetId
end
