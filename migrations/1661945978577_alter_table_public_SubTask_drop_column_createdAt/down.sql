alter table "public"."SubTask" alter column "createdAt" set default now();
alter table "public"."SubTask" alter column "createdAt" drop not null;
alter table "public"."SubTask" add column "createdAt" timestamptz;
