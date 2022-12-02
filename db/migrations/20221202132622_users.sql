-- migrate:up
create table users (
    id string PRIMARY key, 
    name string
)

-- migrate:down
drop table users;
