-- ----------------------------
-- Sequence structure for service_results_id_seq
-- ----------------------------
DROP SEQUENCE IF EXISTS "public"."service_results_id_seq";
CREATE SEQUENCE "public"."service_results_id_seq" 
INCREMENT 1
MINVALUE  1
MAXVALUE 2147483647
START 1
CACHE 1;

-- ----------------------------
-- Table structure for service_results
-- ----------------------------
DROP TABLE IF EXISTS "public"."service_results";
CREATE TABLE "public"."service_results" (
  "id" int4 NOT NULL DEFAULT nextval('service_results_id_seq'::regclass),
  "execution_time" timestamptz(6) NOT NULL,
  "service1_success" bool NOT NULL,
  "service2_success" bool NOT NULL,
  "service1_delay" int8 NOT NULL,
  "service2_delay" int8 NOT NULL
)
;

-- ----------------------------
-- Primary Key structure for table service_results
-- ----------------------------
ALTER TABLE "public"."service_results" ADD CONSTRAINT "service_results_pkey" PRIMARY KEY ("id");
