CREATE TABLE "public"."SubTaskAssignee" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "createdAt" timestamptz NOT NULL DEFAULT now(), "updatedAt" timestamptz NOT NULL DEFAULT now(), "subTaskID" uuid NOT NULL, "userID" uuid NOT NULL, PRIMARY KEY ("id") , FOREIGN KEY ("subTaskID") REFERENCES "public"."SubTask"("id") ON UPDATE restrict ON DELETE cascade, FOREIGN KEY ("userID") REFERENCES "public"."User"("id") ON UPDATE restrict ON DELETE cascade);
CREATE OR REPLACE FUNCTION "public"."set_current_timestamp_updatedAt"()
RETURNS TRIGGER AS $$
DECLARE
  _new record;
BEGIN
  _new := NEW;
  _new."updatedAt" = NOW();
  RETURN _new;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER "set_public_SubTaskAssignee_updatedAt"
BEFORE UPDATE ON "public"."SubTaskAssignee"
FOR EACH ROW
EXECUTE PROCEDURE "public"."set_current_timestamp_updatedAt"();
COMMENT ON TRIGGER "set_public_SubTaskAssignee_updatedAt" ON "public"."SubTaskAssignee" 
IS 'trigger to set value of column "updatedAt" to current timestamp on row update';
CREATE EXTENSION IF NOT EXISTS pgcrypto;
