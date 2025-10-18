CREATE TABLE IF NOT EXISTS  workflows_schedules (
                                                    id  uuid not null primary key ,
                                                    related_to varchar(100) references app_users(clerk_id),
                                                    workflow_type varchar(100) not null ,
                                                    schedule_id varchar(100) not null ,
                                                    created_at timestamp default  now()
)