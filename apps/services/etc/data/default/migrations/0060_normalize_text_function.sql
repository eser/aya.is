-- +goose Up
-- Create a helper function for accent-insensitive text search.
-- Normalizes common accented/special characters to ASCII equivalents.

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION normalize_text(input TEXT) RETURNS TEXT AS $$
BEGIN
  RETURN translate(
    LOWER(input),
    'şığüöçŞİĞÜÖÇéèêëàâäùûüôîïñæøåðþßáíóúýẞ',
    'siguccsigooceeeeaaaeuuuoiinaoaddtssaiouysS'
  );
END;
$$ LANGUAGE plpgsql IMMUTABLE;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION IF EXISTS normalize_text(TEXT);
