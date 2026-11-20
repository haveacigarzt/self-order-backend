CREATE TABLE produk_terkait (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    produk_id_a BIGINT NOT NULL REFERENCES produk(id) ON DELETE CASCADE,
    produk_id_b BIGINT NOT NULL REFERENCES produk(id) ON DELETE CASCADE,

    CONSTRAINT produk_terkait_unique
        UNIQUE (produk_id_a, produk_id_b),

    CONSTRAINT produk_terkait_different
        CHECK (produk_id_a <> produk_id_b)
);

CREATE INDEX idx_produk_terkait_produk_id_a
    ON produk_terkait (produk_id_a);
