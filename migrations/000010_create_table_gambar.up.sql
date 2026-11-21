CREATE TABLE IF NOT EXISTS gambar (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    produk_id BIGINT NOT NULL REFERENCES produk(id) ON DELETE CASCADE,
    file_path VARCHAR(255) NOT NULL,
    width SMALLINT,
    height SMALLINT,
    urutan SMALLINT NOT NULL,
    is_thumbnail BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT gambar_produk_urutan_unique
        UNIQUE (produk_id, urutan),

    CONSTRAINT gambar_urutan_check
        CHECK (urutan BETWEEN 1 AND 5),

    CONSTRAINT gambar_width_check
        CHECK (width IS NULL OR width > 0),

    CONSTRAINT gambar_height_check
        CHECK (height IS NULL OR height > 0)
);

-- Maksimal satu thumbnail untuk setiap produk
CREATE UNIQUE INDEX gambar_satu_thumbnail_per_produk
ON gambar (produk_id)
WHERE is_thumbnail = TRUE;