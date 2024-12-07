package core

import (
	"encoding/json"
	"errors"
	"github.com/sinameshkini/gorawler/internal/repository"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

func (c *Core) Json2Sql(filePath string) (err error) {
	var fields = make(map[string]interface{})

	if filePath == "" {
		return errors.New("Set model json file path.\t gorawler json2sql -m my_model.json")
	}

	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return
	}

	if err = json.Unmarshal(bytes, &fields); err != nil {
		return
	}

	x := strings.Split(filePath, ".")
	tableName := x[len(x)-1]

	if err = repository.Migrator(c.db, tableName, fields); err != nil {
		logrus.Errorln(err)
		return
	}

	if err = repository.Seeder(c.db, tableName, fields); err != nil {
		logrus.Errorln(err)
		return
	}

	return nil
}
