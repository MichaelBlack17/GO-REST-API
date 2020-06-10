CREATE TABLE RequestQueue(
                        Id          bigserial not null primary key,
                        RequestId   bigint not null,
                        ManagerId   bigint,
                        Status      int not null,
                        ValidTime   timestamp
);
