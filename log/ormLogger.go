package log

import (
	"time"
	"fmt"
	"database/sql/driver"
	"reflect"
	"regexp"
	"github.com/jinzhu/gorm"
	"unicode"
)

var sqlRegexp = regexp.MustCompile(`(\$\d+)|\?`)

// Logger default logger
type OrmLogger struct {
	gorm.LogWriter
}

// Print format & print log
func (logger OrmLogger) Print(values ...interface{}) {
	if len(values) > 1 {
		level := values[0]
		source := values[1]
		currentTime := time.Now().Format("2006/01/02 15:04:05")

		d := ""
		var sql string

		if level == "sql" {
			// duration
			d = fmt.Sprintf("%.2fms", float64(values[2].(time.Duration).Nanoseconds() / 1e4) / 100.0);

			// sql
			var formattedValues []string

			for _, value := range values[4].([]interface{}) {
				indirectValue := reflect.Indirect(reflect.ValueOf(value))
				if indirectValue.IsValid() {
					value = indirectValue.Interface()
					if t, ok := value.(time.Time); ok {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", t.Format(time.RFC3339)))
					} else if b, ok := value.([]byte); ok {
						if str := string(b); isPrintable(str) {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", str))
						} else {
							formattedValues = append(formattedValues, "'<binary>'")
						}
					} else if r, ok := value.(driver.Valuer); ok {
						if value, err := r.Value(); err == nil && value != nil {
							formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
						} else {
							formattedValues = append(formattedValues, "NULL")
						}
					} else {
						formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
					}
				} else {
					formattedValues = append(formattedValues, fmt.Sprintf("'%v'", value))
				}
			}

			var formattedValuesLength = len(formattedValues)
			for index, value := range sqlRegexp.Split(values[3].(string), -1) {
				sql += value
				if index < formattedValuesLength {
					sql += formattedValues[index]
				}
			}
		}

		Printf("%s [%s] %v | %s", currentTime, Color.Cyan("ORM"), Color.Blue(fmt.Sprintf("%-6s", "PATH")), source)
		Printf("%s [%s] %v | %-5s | %s", currentTime, Color.Cyan("ORM"), Color.Green(fmt.Sprintf("%-6s", "QUERY")), d, sql)
	}
}

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}
