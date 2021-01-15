create table dbo.BetList
(
    SurebetId    bigint                                     not null,
    SideIndex    tinyint                                    not null,

    BetId        bigint                                     not null,


    Price        decimal(9, 5),
    Stake        decimal(9, 5),
    WinLoss      decimal(9, 5),

    ApiBetId     varchar(1000),
    ApiBetStatus varchar(1000),

    UpdatedAt    datetimeoffset default sysdatetimeoffset() not null,
    CreatedAt    datetimeoffset default sysdatetimeoffset() not null,
    constraint PK_BetList primary key (SurebetId, SideIndex),
)

create type dbo.BetListType as table
(
    SurebetId    bigint  not null,
    SideIndex    tinyint not null,

    BetId        bigint  not null,

    Price        decimal(9, 5),
    Stake        decimal(9, 5),
    WinLoss      decimal(9, 5),

    ApiBetId     varchar(1000),
    ApiBetStatus varchar(1000),

    primary key (SurebetId, SideIndex)
)