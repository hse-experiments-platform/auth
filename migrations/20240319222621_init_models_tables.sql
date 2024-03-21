-- +goose Up
-- +goose StatementBegin
create table metrics
(
    id          bigint primary key generated by default as identity,
    name        text                     not null,
    description text                     not null,
    created_at  timestamp with time zone not null default now(),
    updated_at  timestamp with time zone not null default now()
);

create table problems
(
    id          bigint primary key generated by default as identity,
    name        text                     not null,
    description text                     not null,
    created_at  timestamp with time zone not null default now(),
    updated_at  timestamp with time zone not null default now()
);

create table problem_metrics
(
    problem_id bigint references problems (id) not null,
    metric_id  bigint references metrics (id)  not null
);

create table models
(
    id          bigint primary key generated by default as identity,
    name        text                     not null,
    description text                     not null,
    problem_id  bigint                   not null references problems (id),
    created_at  timestamp with time zone not null default now(),
    updated_at  timestamp with time zone not null default now()
);

create table hyperparameters
(
    id            bigint primary key generated by default as identity,
    name          text                     not null,
    description   text                     not null,
    "type"        text                     not null,
    default_value jsonb                    not null,
    model_id      bigint references models (id),
    created_at    timestamp with time zone not null default now(),
    updated_at    timestamp with time zone not null default now()
);



create table trained_models
(
    id                    bigint primary key generated by default as identity,
    name                  text                     not null,
    description           text                     not null,
    model_id              bigint                   not null references models (id),
    model_training_status model_training_status    not null,
    training_dataset_id   bigint                   not null references datasets (id),
    target_column         text                     not null,
    train_error           text,
    created_at            timestamp with time zone not null default now(),
    updated_at            timestamp with time zone not null default now(),
    launch_id             bigint                   not null
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists metrics;
drop table if exists problems;
drop table if exists problem_metrics;
drop table if exists models;
drop table if exists hyperparameters;
drop table if exists trained_models;
-- +goose StatementEnd
