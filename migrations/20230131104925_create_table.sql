-- +goose Up
-- +goose StatementBegin
create table test(
    id bigint NOT NULL,
    name text,
    meta text,
    status text,
    created_at timestamp,
    updated_at timestamp,
    PRIMARY KEY(id)
);

create index on test(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table test;
-- +goose StatementEnd
