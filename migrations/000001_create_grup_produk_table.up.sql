CREATE TABLE IF NOT EXISTS grup_produk (
  id bigserial PRIMARY KEY,
  nama text NOT NULL,
  created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
  version integer NOT NULL DEFAULT 1
);
