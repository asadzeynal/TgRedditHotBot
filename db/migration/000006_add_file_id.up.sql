ALTER TABLE post_images ADD COLUMN "tg_file_id" BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX ON "tg_file_id" ("post_images");

ALTER TABLE post_videos ADD COLUMN "tg_file_id" BOOLEAN NOT NULL DEFAULT FALSE;

CREATE INDEX ON "tg_file_id" ("post_videos");
