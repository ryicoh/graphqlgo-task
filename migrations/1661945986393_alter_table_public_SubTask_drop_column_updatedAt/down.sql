alter table "public"."SubTask" alter column "updatedAt" set default now();
alter table "public"."SubTask" alter column "updatedAt" drop not null;
alter table "public"."SubTask" add column "updatedAt" timestamptz;
