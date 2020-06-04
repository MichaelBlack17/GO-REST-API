CREATE TABLE Requests(
                      Id        bigserial not null primary key ,
                      UserId    bigint not null,
                      Message   varchar not null,
                      CreateDate timestamp not null
);
