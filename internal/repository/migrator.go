package repository

import (
	"fmt"
	"github.com/sinameshkini/gorawler/internal/utils"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"reflect"
)

func Seeder(db *gorm.DB, table string, fields map[string]interface{}) (err error) {
	var mustInsert = make(map[string]interface{})
	// Iterate over the fields in the JSON model
	for key, value := range fields {
		if value == nil {
			continue
		}
		// Detect field type dynamically
		columnType, _ := detectColumnType(value)
		if columnType == "TABLE" {
			snowfields := make(map[string]interface{})
			if err = utils.JsonAssertion(value, &snowfields); err != nil {
				return err
			}
			if err = Seeder(db, key, snowfields); err != nil {
				return err
			}

			continue
		}

		if columnType == "SLICE" {
			var m = make([]interface{}, 0)

			if err = utils.JsonAssertion(value, &m); err != nil {
				return err
			}

			for _, n := range m {
				snowfields := make(map[string]interface{})
				if err = utils.JsonAssertion(n, &snowfields); err != nil {
					return err
				}
				if err = Seeder(db, key, snowfields); err != nil {
					return err
				}
			}

			continue
		}

		mustInsert[key] = value
	}

	if _, ok := fields["id"]; ok {
		var tmp string
		if _ = db.Table(table).Select("id").Where("id = ?", fields["id"]).Scan(&tmp).Error; tmp == "" {
			if err := db.Table(table).Create(&mustInsert).Error; err != nil {
				return fmt.Errorf("failed to insert record %+v: %w", fields, err)
			}
		}
	} else {
		if err := db.Table(table).Create(&mustInsert).Error; err != nil {
			return fmt.Errorf("failed to insert record %+v: %w", fields, err)
		}
	}
	return nil
}

func Migrator(db *gorm.DB, table string, fields map[string]interface{}) (err error) {
	// Create table if it doesn't exist
	if !db.Migrator().HasTable(table) {
		fmt.Printf("Table %s does not exist. Creating it...\n", table)
		if err := db.Exec(fmt.Sprintf("CREATE TABLE %s (id VARCHAR)", table)).Error; err != nil {
			return fmt.Errorf("failed to create table %s: %w", table, err)
		}
	}

	// Iterate over the fields in the JSON model
	for key, value := range fields {
		if value == nil {
			continue
		}
		// Detect field type dynamically
		columnType, nullable := detectColumnType(value)
		if columnType == "TABLE" {
			snowfields := make(map[string]interface{})
			if err = utils.JsonAssertion(value, &snowfields); err != nil {
				return err
			}
			if err = Migrator(db, key, snowfields); err != nil {
				return err
			}

			continue
		}

		if columnType == "SLICE" {
			var m = make([]interface{}, 0)

			if err = utils.JsonAssertion(value, &m); err != nil {
				return err
			}

			snowfields := make(map[string]interface{})
			if err = utils.JsonAssertion(m[0], &snowfields); err != nil {
				return err
			}
			if err = Migrator(db, key, snowfields); err != nil {
				return err
			}

			continue
		}
		// Check if the column already exists
		if db.Migrator().HasColumn(table, key) {
			if key == "id" && !db.Migrator().HasConstraint(table, fmt.Sprintf("%s_pkey", key)) {
				// Set 'id' as the primary key
				if err := db.Exec(fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT %s_pkey PRIMARY KEY (id)", table, table)).Error; err != nil {
					logrus.Errorln(fmt.Errorf("failed to set 'id' as primary key: %w", err))
				}
			} else {
				fmt.Printf("Column %s already exists in table %s, skipping...\n", key, table)
			}
			continue
		}

		// Attempt to add the column
		fmt.Printf("Adding column %s to table %s...\n", key, table)
		if err := db.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s %s", table, key, columnType, nullable)).Error; err != nil {
			return fmt.Errorf("failed to add column %s: %w", key, err)
		}
	}

	fmt.Printf("Migration for table %s completed successfully!\n", table)
	return nil
}

// detectColumnType maps Go types to SQL types for dynamic migration
func detectColumnType(value interface{}) (string, string) {
	nullable := "NOT NULL"
	if value == nil {
		nullable = "NULL"
	}

	switch reflect.TypeOf(value).Kind() {
	case reflect.String:
		return "VARCHAR(255)", nullable
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "INTEGER", nullable
	case reflect.Float32, reflect.Float64:
		return "FLOAT", nullable
	case reflect.Bool:
		return "BOOLEAN", nullable
	case reflect.Map:
		return "TABLE", nullable
	case reflect.Slice:
		return "SLICE", nullable
	default:
		return "TEXT", nullable
	}
}
