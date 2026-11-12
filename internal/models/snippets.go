package models

import (
	"database/sql"
	"errors"
	"time"
)

// Define a snippet type to hold the data for an individual snippet
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a snippet model type which a sql.DB connectio pool
type SnippetModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the database
func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {

	// Write the sql statement we want to execute
	stmt := `INSERT INTO snippets (title, content, created, expires)
			VALUES (?, ?, datetime('now'), datetime('now', '+' || ? || ' days'))`

	//stmt := `INSERT INTO snippets (title, content, created, expires)
	//		VALUES ($1, $2, NOW() AT TIME ZONE 'UTC', (NOW() AT TIME ZONE 'UTC') + ($3 * INTERVAL '1 day'))`

	result, err := m.DB.Exec(stmt, title, content, expires)

	if err != nil {
		return 0, err
	}

	// This method on the result to get the id of our newly inserted record in the snippets table
	id, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	//
	return int(id), nil
}

// This will return a specif snippet based on its id
func (m *SnippetModel) Get(id int) (Snippet, error) {

	// Initialize a new zeroed Snippet struct
	var s Snippet

	stmt := `SELECT id, title, content, created, expires FROM snippets
			WHERE expires > datetime('now') AND id = ?`

	//stmt := `SELECT id, title, content, created, expires FROM snippets
	//		WHERE expires > (NOW() AT TIME ZONE 'UTC') AND id = $1`

	// Scan copy the values from each field in the snippet struct
	err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Snippet{}, ErrNoRecord
		} else {
			return Snippet{}, err
		}
	}

	return s, nil
}

// This will return the 10 most recently created snippets
func (m *SnippetModel) Latest() ([]Snippet, error) {

	stmt := `SELECT id, title, content, created, expires FROM snippets
			WHERE expires > datetime('now') ORDER BY id DESC LIMIT 10`

	//stmt := `SELECT id, title, content, created, expires FROM snippets
	//		WHERE expires > NOW() ORDER BY id DESC LIMIT 10`

	// Use to Query method on the connection pool to execute our sql statement
	rows, err := m.DB.Query(stmt)

	if err != nil {
		return nil, err
	}

	// To ensure the sql.Rows resultset is always porperly closed before
	// the Latest() method returns
	defer rows.Close()

	// Initialize an empty slice to hold the snippet structs
	var snippets []Snippet

	for rows.Next() {
		var s Snippet

		err := rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// If everything went ok then return the snippets slice
	return snippets, nil
}

func (m *SnippetModel) Update(title string, content string, expires int, id int) error {

	stmt := "UPDATE snippets SET title = ?, content = ?, expires = DATE_ADD(created, INTERVAL ? DAY) WHERE id = ?"

	//stmt := `UPDATE snippets SET title = $1, content = $2, expires = created + ($3 * INTERVAL '1 day')
	//		WHERE id = $4`
	_, err := m.DB.Exec(stmt, title, content, expires, id)

	return err
}

func (m *SnippetModel) Delete(id int) error {

	stmt := "DELETE FROM snippets WHERE id = ?"
	//stmt := `DELETE FROM snippets WHERE id = $1`
	_, err := m.DB.Exec(stmt, id)

	return err
}
