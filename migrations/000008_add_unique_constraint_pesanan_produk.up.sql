ALTER TABLE pesanan_produk
ADD CONSTRAINT pesanan_produk_unique
UNIQUE (pesanan_id, produk_id);