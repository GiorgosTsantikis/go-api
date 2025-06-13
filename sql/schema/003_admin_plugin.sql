-- +goose Up
alter table "user" add column "role" text;

alter table "user" add column "banned" boolean;

alter table "user" add column "banReason" text;

alter table "user" add column "banExpires" timestamp;

alter table "session" add column "impersonatedBy" text;

-- +goose Down
alter table "session" drop column "impersonatedBy";
alter table "user" drop column "banExpires";
alter table "user" drop column "banReason";
alter table "user" drop column "banned";
alter table "user" drop column "role";