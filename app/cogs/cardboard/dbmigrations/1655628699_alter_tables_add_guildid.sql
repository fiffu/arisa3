ALTER TABLE "aliases"     ADD COLUMN guildid TEXT;
ALTER TABLE "tag_promote" ADD COLUMN guildid TEXT;
ALTER TABLE "tag_demote"  ADD COLUMN guildid TEXT;
ALTER TABLE "tag_omit"    ADD COLUMN guildid TEXT;
