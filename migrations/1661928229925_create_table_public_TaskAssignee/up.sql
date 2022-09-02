CREATE TABLE "public"."TaskAssignee" ("id" uuid NOT NULL DEFAULT gen_random_uuid(), "created_at" timestamptz NOT NULL DEFAULT now(), "updated_at" timestamptz NOT NULL DEFAULT now(), "taskID" uuid NOT NULL, "userID" uuid NOT NULL, PRIMARY KEY ("id") , FOREIGN KEY ("taskID") REFERENCES "public"."Task"("id") ON UPDATE restrict ON DELETE cascade, FOREIGN KEY ("userID") REFERENCES "public"."User"("id") ON UPDATE restrict ON DELETE cascade);
CREATE OR REPLACE FUNCTION "public"."set_current_timestamp_updated_at"()
RETURNS TRIGGER AS $$
DECLARE
  _new record;
BEGIN
  _new := NEW;
  _new."updated_at" = NOW();
  RETURN _new;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER "set_public_TaskAssignee_updated_at"
BEFORE UPDATE ON "public"."TaskAssignee"
FOR EACH ROW
EXECUTE PROCEDURE "public"."set_current_timestamp_updated_at"();
COMMENT ON TRIGGER "set_public_TaskAssignee_updated_at" ON "public"."TaskAssignee" 
IS 'trigger to set value of column "updated_at" to current timestamp on row update';
CREATE EXTENSION IF NOT EXISTS pgcrypto;
