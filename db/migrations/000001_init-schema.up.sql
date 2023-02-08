CREATE TABLE "posts" (
                         "id" varchar PRIMARY KEY,
                         "title" varchar NOT NULL,
                         "url" varchar NOT NULL,
                         "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "post_images" (
                               "id" bigserial PRIMARY KEY,
                               "post" varchar NOT NULL,
                               "url" varchar NOT NULL
);

CREATE TABLE "post_videos" (
                               "id" bigserial PRIMARY KEY,
                               "post" varchar NOT NULL,
                               "height" int NOT NULL,
                               "width" int NOT NULL,
                               "duration" int NOT NULL,
                               "url" varchar NOT NULL
);

CREATE INDEX ON "post_images" ("post");

CREATE INDEX ON "post_videos" ("post");

ALTER TABLE "post_images" ADD FOREIGN KEY ("post") REFERENCES "posts" ("id");

ALTER TABLE "post_videos" ADD FOREIGN KEY ("post") REFERENCES "posts" ("id");
