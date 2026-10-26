package errorIgnore

import "database/sql"

var (
	ignoredErrors = map[error]bool{
		sql.ErrNoRows: true,
	}
)

func Ignore(err error) error {
	if _, ok := ignoredErrors[err]; ok {
		return nil
	}
	return err
}
