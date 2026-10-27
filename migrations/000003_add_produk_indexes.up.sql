CREATE INDEX IF NOT EXISTS produk_nama_idx ON produk USING GIN (to_tsvector('simple', nama));
