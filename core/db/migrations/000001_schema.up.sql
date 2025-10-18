CREATE  TABLE  IF NOT EXISTS  app_users (
                                            id pg_catalog.uuid not null  primary key,
                                            clerk_id varchar(100) not null,
                                            createdAt timestamp default  now()
)