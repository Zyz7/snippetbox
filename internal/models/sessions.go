
package models

import (
    "database/sql"
    "time"
)

type SessionStore struct {
    DB *sql.DB
}

func (s *SessionStore) Find(token string) ([]byte, bool, error) {
    var (
        data   []byte
        expiry int64
    )

    stmt := `SELECT data, expiry
             FROM sessions
             WHERE token = ?`

    err := s.DB.QueryRow(stmt, token).Scan(&data, &expiry)

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, false, nil
        }
        return nil, false, err
    }

    // Check whether the session has expired.
    if time.Now().Unix() > expiry {
        return nil, false, nil
    }

    return data, true, nil
}

func (s *SessionStore) Commit(token string, data []byte, expiry time.Time) error {
    stmt := `INSERT INTO sessions (token, data, expiry)
             VALUES (?, ?, ?)
             ON CONFLICT(token) DO UPDATE SET
                 data = excluded.data,
                 expiry = excluded.expiry`

    _, err := s.DB.Exec(stmt, token, data, expiry.Unix())
    return err
}

func (s *SessionStore) Delete(token string) error {
    stmt := `DELETE FROM sessions WHERE token = ?`

    _, err := s.DB.Exec(stmt, token)
    return err
}
