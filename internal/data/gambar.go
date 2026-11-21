package data

import (
	"database/sql"
	"time"
)

type Gambar struct {
	ID          int64         `json:"id"`
	ProdukID    int64         `json:"produk_id"`
	FilePath    string        `json:"file_path"`
	Width       sql.NullInt16 `json:"width"`
	Height      sql.NullInt16 `json:"height"`
	Urutan      int16         `json:"urutan"`
	IsThumbnail bool          `json:"is_thumbnail"`
	CreatedAt   time.Time     `json:"created_at"`
}

type GambarModel struct {
	DB *sql.DB
}

func (m GambarModel) GetByProdukIDAndUrutan(produkID int64, urutan int16) (string, error) {
	query := `
			SELECT file_path
			FROM gambar
			WHERE produk_id = $1 AND urutan = $2
    `

	var filePath string

	err := m.DB.QueryRow(query, produkID, urutan).Scan(&filePath)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func (m GambarModel) GetAllByProdukID(produkID int64) ([]Gambar, error) {
	query := `
		SELECT file_path, urutan
		FROM gambar
		WHERE produk_id = $1
		ORDER BY urutan ASC
	`

	rows, err := m.DB.Query(query, produkID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	gambarList := []Gambar{}

	for rows.Next() {
		var gambar Gambar

		err := rows.Scan(
			&gambar.FilePath,
			&gambar.Urutan,
		)
		if err != nil {
			return nil, err
		}

		gambarList = append(gambarList, gambar)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return gambarList, nil
}

func (m GambarModel) Upsert(gambar *Gambar) error {
	query := `
        INSERT INTO gambar (
            produk_id,
            file_path,
            width,
            height,
            urutan,
            is_thumbnail
        )
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (produk_id, urutan)
        DO UPDATE SET
            file_path = EXCLUDED.file_path,
            width = EXCLUDED.width,
            height = EXCLUDED.height,
            is_thumbnail = EXCLUDED.is_thumbnail
    `

	args := []any{
		gambar.ProdukID,
		gambar.FilePath,
		gambar.Width,
		gambar.Height,
		gambar.Urutan,
		gambar.IsThumbnail,
	}

	_, err := m.DB.Exec(query, args...)
	return err
}
