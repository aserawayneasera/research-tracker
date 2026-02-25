package store

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/aserawayneasera/research-tracker/internal/models"
)

type SQLitePaperStore struct {
	db *sql.DB
}

func NewSQLitePaperStore(db *sql.DB) *SQLitePaperStore {
	return &SQLitePaperStore{db: db}
}

var ErrNotFound = errors.New("not found")

func (s *SQLitePaperStore) Create(p *models.Paper) (*models.Paper, error) {
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now

	if err := p.Validate(); err != nil {
		return nil, err
	}

	res, err := s.db.Exec(
		`INSERT INTO papers(title, venue, year, status, tags, notes, created_at, updated_at)
		 VALUES(?, ?, ?, ?, ?, ?, ?, ?)`,
		p.Title, p.Venue, p.Year, p.Status, p.Tags, p.Notes, p.CreatedAt.Format(time.RFC3339), p.UpdatedAt.Format(time.RFC3339),
	)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	return s.Get(id)
}

func (s *SQLitePaperStore) Get(id int64) (*models.Paper, error) {
	row := s.db.QueryRow(`SELECT id, title, venue, year, status, tags, notes, created_at, updated_at FROM papers WHERE id = ?`, id)
	var p models.Paper
	var created, updated string
	if err := row.Scan(&p.ID, &p.Title, &p.Venue, &p.Year, &p.Status, &p.Tags, &p.Notes, &created, &updated); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	p.CreatedAt, _ = time.Parse(time.RFC3339, created)
	p.UpdatedAt, _ = time.Parse(time.RFC3339, updated)
	return &p, nil
}

func (s *SQLitePaperStore) Update(id int64, p *models.Paper) (*models.Paper, error) {
	existing, err := s.Get(id)
	if err != nil {
		return nil, err
	}

	// apply fields
	existing.Title = p.Title
	existing.Venue = p.Venue
	existing.Year = p.Year
	existing.Status = p.Status
	existing.Tags = p.Tags
	existing.Notes = p.Notes
	existing.UpdatedAt = time.Now().UTC()

	if err := existing.Validate(); err != nil {
		return nil, err
	}

	_, err = s.db.Exec(
		`UPDATE papers
		 SET title=?, venue=?, year=?, status=?, tags=?, notes=?, updated_at=?
		 WHERE id=?`,
		existing.Title, existing.Venue, existing.Year, existing.Status, existing.Tags, existing.Notes, existing.UpdatedAt.Format(time.RFC3339), id,
	)
	if err != nil {
		return nil, err
	}
	return s.Get(id)
}

func (s *SQLitePaperStore) Delete(id int64) error {
	res, err := s.db.Exec(`DELETE FROM papers WHERE id=?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *SQLitePaperStore) List(params ListParams) ([]models.Paper, int, error) {
	limit := params.Limit
	offset := params.Offset
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}

	where := []string{"1=1"}
	args := []any{}

	if strings.TrimSpace(params.Status) != "" {
		where = append(where, "status = ?")
		args = append(args, strings.ToLower(strings.TrimSpace(params.Status)))
	}
	if params.Year != 0 {
		where = append(where, "year = ?")
		args = append(args, params.Year)
	}
	if strings.TrimSpace(params.Q) != "" {
		where = append(where, "(title LIKE ? OR venue LIKE ? OR tags LIKE ?)")
		q := "%" + strings.TrimSpace(params.Q) + "%"
		args = append(args, q, q, q)
	}

	whereSQL := strings.Join(where, " AND ")

	// total count
	var total int
	if err := s.db.QueryRow("SELECT COUNT(*) FROM papers WHERE "+whereSQL, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// list rows
	argsList := append([]any{}, args...)
	argsList = append(argsList, limit, offset)

	rows, err := s.db.Query(
		`SELECT id, title, venue, year, status, tags, notes, created_at, updated_at
		 FROM papers
		 WHERE `+whereSQL+`
		 ORDER BY updated_at DESC
		 LIMIT ? OFFSET ?`,
		argsList...,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	out := []models.Paper{}
	for rows.Next() {
		var p models.Paper
		var created, updated string
		if err := rows.Scan(&p.ID, &p.Title, &p.Venue, &p.Year, &p.Status, &p.Tags, &p.Notes, &created, &updated); err != nil {
			return nil, 0, err
		}
		p.CreatedAt, _ = time.Parse(time.RFC3339, created)
		p.UpdatedAt, _ = time.Parse(time.RFC3339, updated)
		out = append(out, p)
	}
	return out, total, rows.Err()
}
