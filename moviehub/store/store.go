package store

import (
	"errors"
"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

type Store struct {
	Pool *pgxpool.Pool
}

func (s *Store) ListMembers(ctx context.Context) ([]Member, error) {
    rows, err := s.Pool.Query(ctx,
        `SELECT id, email, display_name, role FROM members ORDER BY id`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    out := []Member{}
    for rows.Next() {
        var m Member
        if err := rows.Scan(&m.ID, &m.Email, &m.DisplayName, &m.Role); err != nil {
            return nil, err
        }
        out = append(out, m)
    }
    return out, rows.Err()
}