CREATE  TABLE IF NOT EXISTS  kainos_user (
                                             id uuid primary key ,
                                             clerk_id varchar not null ,
                                             first_name varchar ,
                                             email varchar not null  unique,
                                             created_at timestamp not null  default  now(),
                                             deleted_at timestamp,
                                             updated_at timestamp
    );

CREATE  TABLE  IF NOT EXISTS  kainos_workflow (
                                                  id uuid primary key,
                                                  workflow_name varchar not null ,
                                                  workflow_description text not null,
                                                  created_at timestamp not null  default  now(),
    deleted_at timestamp,
    updated_at timestamp,
    price pg_catalog.float8 default 0.0
    );

CREATE TABLE  IF NOT EXISTS  kainos_user_workflow (
                                                      id uuid primary key ,
                                                      workflow_id uuid not null references kainos_workflow(id),
    customer_id pg_catalog.uuid not null  references  kainos_user(id),
    meta_data pg_catalog.jsonb default '{}',
    cron_time varchar,
    status varchar default  'OFF',
    created_at timestamp default  now(),
    updated_at timestamp
    );

CREATE  TABLE  IF NOT EXISTS  system_defined_analysis (
                                                          id uuid primary key,
                                                          analysis_type varchar,
                                                          description varchar
);
CREATE  TABLE  IF NOT EXISTS  kainos_user_analysis (
                                                       id uuid primary key,
                                                       description varchar, -- ai generated
                                                       s3_url varchar,
                                                       customer_id uuid not null references kainos_user(id),
    created_at timestamp default  now()
    );
