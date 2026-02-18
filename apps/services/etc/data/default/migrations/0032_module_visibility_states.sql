-- Rename hide_* columns to feature_* and convert from BOOLEAN to TEXT
-- with three states: 'public', 'hidden', 'disabled'.
-- Backward-compatible: false → 'public', true → 'hidden'.

-- feature_relations (was hide_relations)
ALTER TABLE "profile" RENAME COLUMN "hide_relations" TO "feature_relations";
ALTER TABLE "profile"
  ALTER COLUMN "feature_relations" DROP DEFAULT,
  ALTER COLUMN "feature_relations" DROP NOT NULL,
  ALTER COLUMN "feature_relations" TYPE TEXT USING CASE WHEN feature_relations::BOOLEAN THEN 'hidden' ELSE 'public' END;
ALTER TABLE "profile"
  ALTER COLUMN "feature_relations" SET DEFAULT 'public',
  ALTER COLUMN "feature_relations" SET NOT NULL;

-- feature_links (was hide_links)
ALTER TABLE "profile" RENAME COLUMN "hide_links" TO "feature_links";
ALTER TABLE "profile"
  ALTER COLUMN "feature_links" DROP DEFAULT,
  ALTER COLUMN "feature_links" DROP NOT NULL,
  ALTER COLUMN "feature_links" TYPE TEXT USING CASE WHEN feature_links::BOOLEAN THEN 'hidden' ELSE 'public' END;
ALTER TABLE "profile"
  ALTER COLUMN "feature_links" SET DEFAULT 'public',
  ALTER COLUMN "feature_links" SET NOT NULL;

-- feature_qa (was hide_qa)
ALTER TABLE "profile" RENAME COLUMN "hide_qa" TO "feature_qa";
ALTER TABLE "profile"
  ALTER COLUMN "feature_qa" DROP DEFAULT,
  ALTER COLUMN "feature_qa" DROP NOT NULL,
  ALTER COLUMN "feature_qa" TYPE TEXT USING CASE WHEN feature_qa::BOOLEAN THEN 'hidden' ELSE 'public' END;
ALTER TABLE "profile"
  ALTER COLUMN "feature_qa" SET DEFAULT 'public',
  ALTER COLUMN "feature_qa" SET NOT NULL;
