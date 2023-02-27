ALTER TABLE post_images ADD COLUMN "tg_file_id" varchar NOT NULL DEFAULT '';

CREATE INDEX ON "post_images" ("tg_file_id");

ALTER TABLE post_videos ADD COLUMN "tg_file_id" varchar NOT NULL DEFAULT '';

CREATE INDEX ON "post_videos" ("tg_file_id");
