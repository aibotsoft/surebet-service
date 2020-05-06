create table dbo.Calc
(
    Profit          decimal(9, 5),
    FirstName       varchar(1000),
    SecondName      varchar(1000),
    LowerWinIndex   tinyint,
    HigherWinIndex  tinyint,
    FirstIndex      tinyint,
    SecondIndex     tinyint,
    WinDiff         decimal(9, 5),
    WinDiffRel      decimal(9, 5),
    FortedSurebetId int                                        not null,
    SurebetId       bigint                                     not null,

    CreatedAt       datetimeoffset default sysdatetimeoffset() not null,
    constraint PK_Calc primary key (SurebetId),
)