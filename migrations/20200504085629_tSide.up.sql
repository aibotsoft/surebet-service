create table dbo.Side
(
    SurebetId     bigint                                     not null,
    SideIndex     tinyint                                        not null,

    ServiceName   varchar(1000),
    SportName     varchar(1000),
    LeagueName    varchar(1000),
    Home          varchar(1000),
    Away          varchar(1000),
    MarketName    varchar(1000),
    Price         decimal(9, 5),
    Initiator     bit,
    Starts        datetimeoffset,
    EventId       int,

    CheckId       bigint,
    AccountId     tinyint,
    AccountLogin  varchar(1000),
    CheckStatus   varchar(1000),
    StatusInfo    varchar(1000),
    CountLine     int,
    CountEvent    int,
    AmountEvent   int,
    MinBet        decimal(9, 5),
    MaxBet        decimal(9, 5),
    Balance       int,
    CheckPrice    decimal(9, 5),
    Currency      decimal(9, 5),
    CheckDone     bigint,

    CalcStatus    varchar(1000),
    MaxStake      decimal(9, 5),
    MinStake      decimal(9, 5),
    MaxWin        decimal(9, 5),
    Stake         decimal(9, 5),
    Win           decimal(9, 5),
    IsFirst       bit,

    ToBetId       bigint,
    TryCount      int,

    BetStatus     varchar(1000),
    BetStatusInfo varchar(1000),
    Start         bigint,
    Done          bigint,
    BetPrice      decimal(9, 5),
    BetStake      decimal(9, 5),
    ApiBetId      varchar(1000),

    UpdatedAt     datetimeoffset default sysdatetimeoffset() not null,
    CreatedAt     datetimeoffset default sysdatetimeoffset() not null,
    constraint PK_Side primary key (SurebetId, SideIndex),
)