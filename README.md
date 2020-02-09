# Start a Postgres database container

docker run --name postgres -e POSTGRES_PASSWORD=password -d -p 5432:5432  postgres

# Create the events table

create table events
(
    id          serial primary key,
    title       varchar(50)  not null,
    description varchar(200) not null
);
