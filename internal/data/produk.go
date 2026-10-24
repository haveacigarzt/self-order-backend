package data

import (
	"time"
)

type Produk struct {
	ID        int64     "json:\"id\""
	GroupID   int64     "json:\"group_id\""
	Nama      string    "json:\"nama\""
	Slug      string    "json:\"slug\""
	Deskripsi string    "json:\"deskripsi\""
	Harga     uint32    "json:\"harga\""
	Stok      uint16    "json:\"stok\""
	CreatedAt time.Time "json:\"-\""
	UpdatedAt time.Time "json:\"-\""
	Version   int32     "json:\"version\""
}
