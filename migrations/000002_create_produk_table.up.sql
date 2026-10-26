CREATE TABLE IF NOT EXISTS produk (
  id bigserial PRIMARY KEY,
  grup_id text NOT NULL,
  nama text NOT NULL,
  slug text NOT NULL,
  deskripsi text NOT NULL,
  harga BIGINT NOT NULL CHECK (harga BETWEEN 0 AND 4294967295),
  stok INTEGER NOT NULL CHECK (stok BETWEEN 0 AND 65535),
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  version integer NOT NULL DEFAULT 1
);
