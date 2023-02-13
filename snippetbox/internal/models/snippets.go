package models

import (
	"database/sql"
	"errors"
	"time"
)

// Snippet to hold the data for individual records of snippet type.
type Snippet struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Snippet model which wraps the sql.DB connection pool.
type SnippetModel struct {
	DB *sql.DB
}

// func (m *SnippetModel) Insert(title string, content string, expires int) (int, error) {
func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	query := `INSERT INTO snippets (title, content, created, expires) 
	VALUES ($1, $2, NOW(), NOW() + $3) 
	RETURNING id;`

	// Below is for MySQL
	// result, err := m.DB.Exec(query, title, content, expires)
	// if err != nil {
	// 	return 0, err
	// }
	// defer result.Close()

	// id, err := result.LastInsertId()
	// if err != nil {
	// 	return 0, err
	// }
	// return int(id), nil
	var id int
	err := m.DB.QueryRow(query, title, content, expires).Scan(&id)
	// fmt.Println("New record ID is:", id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *SnippetModel) Get(id int) (*Snippet, error) {
	query := `SELECT id, title, content, created, expires FROM snippets WHERE expires > NOW() AND id = $1`

	row := m.DB.QueryRow(query, id)

	snippet := &Snippet{}

	err := row.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return snippet, nil
}

func (m *SnippetModel) Latest() ([]*Snippet, error) {
	stmt := `SELECT id, title, content, created, expires FROM snippets WHERE expires > NOW() ORDER BY id DESC LIMIT 10;`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*Snippet{}

	for rows.Next() {
		snippet := &Snippet{}
		err = rows.Scan(&snippet.ID, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expires)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, snippet)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
