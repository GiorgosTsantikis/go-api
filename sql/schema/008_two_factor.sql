-- +goose Up
alter table "user" add column "twoFactorEnabled" boolean;

create table "twoFactor" ("id" text not null primary key, "secret" text not null, "backupCodes" text not null, "userId" text not null references "user" ("id"));

-- +goose Down
drop table "twoFactor";
alter table "user" drop column "twoFactorEnabled";