alter table "public"."Project"
  add constraint "Project_visibility_fkey"
  foreign key ("visibility")
  references "public"."Visibility"
  ("value") on update restrict on delete restrict;
