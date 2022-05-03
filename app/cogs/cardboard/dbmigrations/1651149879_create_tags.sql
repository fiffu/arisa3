CREATE TABLE "tag_promote" (
    tag TEXT NOT NULL UNIQUE
);
CREATE TABLE "tag_demote" (
    tag TEXT NOT NULL UNIQUE
);
CREATE TABLE "tag_omit" (
    tag TEXT NOT NULL UNIQUE
);
