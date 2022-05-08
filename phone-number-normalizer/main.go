package main

import (
	"database/sql"
	"fmt"
	"regexp"

	_ "github.com/lib/pq"
)

func normalize(phone string) string {
	re := regexp.MustCompile("\\D")
	return re.ReplaceAllString(phone, "")
}

const (
	user     = "docker"
	password = "docker"
	dbname   = "db_phone"
	port     = 5432
)

type phone struct {
	id     int
	number string
}

func main() {
	connStr := fmt.Sprintf(`
	 user=%s
	 dbname=%s
	 password=%s
	 port=%d
	 sslmode=disable`,
		user, dbname, password, port)
	db, err := sql.Open("postgres", connStr)
	defer db.Close()
	must(err)
	must(db.Ping())
	createTable(db, "phone_numbers")
	seed(db)
	normalizeDb(db)
}

func normalizeDb(db *sql.DB) {
	phones, err := getPhones(db)
	must(err)
	for _, phone := range phones {
		fmt.Printf("Working on ...%+v\n", phone)
		number := normalize(phone.number)
		if number != phone.number {
			fmt.Printf("Updating or removing...%+v\n", number)
			existing, err := findPhone(db, number)
			must(err)
			if existing != nil {
				must(deletePhone(db, phone.id))
			} else {
				phone.number = number
				must(updatePhone(db, phone))
			}
		} else {
			fmt.Println("No need to update!!!")
		}
	}
}

func createTable(db *sql.DB, name string) {
	statement := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s(
		id SERIAL,
		value VARCHAR(255)
	)`, name)
	_, err := db.Exec(statement)
	must(err)
}

func insert(db *sql.DB, phone string) (int, error) {
	var id int
	statement := "INSERT INTO phone_numbers(value) VALUES($1) RETURNING id"
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func getPhone(db *sql.DB, id int) (string, error) {
	var number string
	row := db.QueryRow("SELECT * FROM phone_numbers WHERE id = $1", id)
	err := row.Scan(&id, &number)
	if err != nil {
		return "", err
	}
	return number, nil
}

func findPhone(db *sql.DB, p string) (*phone, error) {
	var number string
	var id int
	row := db.QueryRow("SELECT * FROM phone_numbers WHERE value = $1", p)
	err := row.Scan(&id, &number)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &phone{
		id:     id,
		number: number,
	}, nil
}

func deletePhone(db *sql.DB, id int) error {
	statement := "DELETE FROM phone_numbers WHERE id = $1"
	_, err := db.Exec(statement, id)
	return err
}

func updatePhone(db *sql.DB, p phone) error {
	statement := "UPDATE phone_numbers SET value = $2 WHERE id = $1"
	_, err := db.Exec(statement, p.id, p.number)
	return err
}

func getPhones(db *sql.DB) ([]phone, error) {
	var p phone
	var ret []phone
	rows, err := db.Query("SELECT * FROM phone_numbers")
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if err := rows.Scan(&p.id, &p.number); err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return ret, nil
}

func seed(db *sql.DB) {
	_, err := insert(db, "1234567890")
	must(err)
	_, err = insert(db, "123 456 7891")
	must(err)
	_, err = insert(db, "(123) 456 7892")
	must(err)
	_, err = insert(db, "(123) 456-7893")
	must(err)
	_, err = insert(db, "123-456-7894")
	must(err)
	_, err = insert(db, "123-456-7890")
	must(err)
	_, err = insert(db, "1234567892")
	must(err)
	_, err = insert(db, "(123)456-7892")
	must(err)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// func isDigit(n rune) bool {
// 	return n >= '0' && n <= '9'
// }

// func normalize(phone string) string {
// 	var buf bytes.Buffer
// 	for _, c := range phone {
// 		if isDigit(c) {
// 			buf.WriteRune(c)
// 		}
// 	}
// 	return buf.String()
// }
