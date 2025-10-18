CREATE  TABLE  IF NOT EXISTS  app_users (
    id pg_catalog.uuid not null  primary key,
    clerk_id varchar(100) not null,
    createdAt timestamp default  now()
);
CREATE TABLE IF NOT EXISTS  workflows_schedules (
    id  pg_catalog.uuid not null primary key ,
    related_to varchar(100) references app_users(clerk_id),
    workflow_type varchar(100) not null ,
    schedule_id varchar(100) not null ,
    created_at timestamp default  now()
);