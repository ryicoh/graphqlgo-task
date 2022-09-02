alter table "public"."Project" add column "visibility" text
 not null default 'ProjectMemberOnly';
