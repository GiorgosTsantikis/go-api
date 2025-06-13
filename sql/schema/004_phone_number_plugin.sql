-- +goose Up
alter table "user" add column "phoneNumber" text unique;

alter table "user" add column "phoneNumberVerified" boolean;

-- +goose Down
alter table "user" drop column "phoneNumberVerified";
alter table "user" drop column "phoneNumber";