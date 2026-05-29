package storage

import (
	"database/sql"
	"time"

	"github.com/adrg/xdg"
	_ "modernc.org/sqlite"
)

type Entry struct {
	Word        string
	Definition  string
	Translation string
	LookupCount int
	LastLookedUp time.Time
}

var db *sql.DB

func Init() error {
	path, err := xdg.DataFile("vtdict/history.db")
	if err != nil {
		return err
	}

	db, err = sql.Open("sqlite", path)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS history (
		word TEXT PRIMARY KEY,
		definition TEXT,
		translation TEXT,
		lookup_count INTEGER DEFAULT 1,
		last_looked_up DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	return err
}

func Save(word, definition, translation string) error {
	_, err := db.Exec(`INSERT INTO history (word, definition, translation, lookup_count, last_looked_up)
		VALUES (?, ?, ?, 1, CURRENT_TIMESTAMP)
		ON CONFLICT(word) DO UPDATE SET
			definition = excluded.definition,
			translation = excluded.translation,
			lookup_count = lookup_count + 1,
			last_looked_up = CURRENT_TIMESTAMP`,
		word, definition, translation)
	return err
}

func Get(word string) (*Entry, error) {
	var e Entry
	var ts string
	err := db.QueryRow(`SELECT word, definition, translation, lookup_count,
		datetime(last_looked_up, 'localtime') FROM history WHERE word = ?`, word).
		Scan(&e.Word, &e.Definition, &e.Translation, &e.LookupCount, &ts)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	e.LastLookedUp, _ = time.Parse("2006-01-02 15:04:05", ts)
	return &e, nil
}

func GetHistory(limit int) ([]Entry, error) {
	rows, err := db.Query(`SELECT word, definition, translation, lookup_count,
		datetime(last_looked_up, 'localtime') FROM history ORDER BY last_looked_up DESC LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var e Entry
		var ts string
		if err := rows.Scan(&e.Word, &e.Definition, &e.Translation, &e.LookupCount, &ts); err != nil {
			continue
		}
		e.LastLookedUp, _ = time.Parse("2006-01-02 15:04:05", ts)
		entries = append(entries, e)
	}
	return entries, nil
}

func Close() {
	if db != nil {
		db.Close()
	}
}

func ClearHistory() error {
	_, err := db.Exec("DELETE FROM history")
	return err
}

func ClearBefore(d time.Duration) (int64, error) {
	cutoff := time.Now().UTC().Add(-d)
	res, err := db.Exec("DELETE FROM history WHERE last_looked_up < ?", cutoff.Format("2006-01-02 15:04:05"))
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func ClearWord(word string) (bool, error) {
	res, err := db.Exec("DELETE FROM history WHERE word = ?", word)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

func ClearWords(words []string) (int64, error) {
	var total int64
	for _, w := range words {
		res, err := db.Exec("DELETE FROM history WHERE word = ?", w)
		if err != nil {
			return total, err
		}
		n, _ := res.RowsAffected()
		total += n
	}
	return total, nil
}
