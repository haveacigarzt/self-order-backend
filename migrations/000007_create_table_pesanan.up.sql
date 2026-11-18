
-- create_table_pesanan.up.sql

CREATE TABLE pesanan (
    id BIGSERIAL PRIMARY KEY,
    no_meja VARCHAR(20) NOT NULL,

    status VARCHAR(30) NOT NULL DEFAULT 'pending_payment',

    id_transaksi VARCHAR(100) UNIQUE,

    total_amount NUMERIC(12,2) NOT NULL DEFAULT 0
        CHECK (total_amount >= 0),

    nama_pemesan VARCHAR(100),

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    catatan TEXT,

    CONSTRAINT pesanan_status_check
    CHECK (
        status IN (
            'pending_payment',
            'paid',
            'cooking',
            'served',
            'confirmed'
        )
    )
);

CREATE TABLE pesanan_produk (
    id BIGSERIAL PRIMARY KEY,

    pesanan_id BIGINT NOT NULL
        REFERENCES pesanan(id) ON DELETE RESTRICT,

    produk_id BIGINT NOT NULL
        REFERENCES produk(id) ON DELETE RESTRICT,

    qty INTEGER NOT NULL CHECK (qty > 0),

    catatan TEXT,

    harga_saat_pesan NUMERIC(12,2) NOT NULL
        CHECK (harga_saat_pesan >= 0),

    subtotal NUMERIC(12,2)
        GENERATED ALWAYS AS (
            qty * harga_saat_pesan
        ) STORED
);

CREATE TABLE status_timeline (
    id BIGSERIAL PRIMARY KEY,

    pesanan_id BIGINT NOT NULL
        REFERENCES pesanan(id) ON DELETE RESTRICT,

    status VARCHAR(30) NOT NULL,

    actor_type VARCHAR(20) NOT NULL DEFAULT 'system',

    actor_id BIGINT,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT timeline_status_check
    CHECK (
        status IN (
            'pending_payment',
            'paid',
            'cooking',
            'served',
            'confirmed'
        )
    ),

    CONSTRAINT timeline_actor_check
    CHECK (
        actor_type IN (
            'system',
            'customer',
            'cashier',
            'kitchen'
        )
    )
);

CREATE INDEX idx_pesanan_status
ON pesanan(status);

CREATE INDEX idx_pesanan_created_at
ON pesanan(created_at DESC);

CREATE INDEX idx_pesanan_produk_pesanan
ON pesanan_produk(pesanan_id);

CREATE INDEX idx_status_timeline_pesanan
ON status_timeline(pesanan_id, created_at);
