package cmd

import (
	"github.com/sinameshkini/gorawler/internal/core"
	"github.com/sinameshkini/gorawler/internal/repository"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(json2sqlCmd)

	json2sqlCmd.Flags().BoolVarP(&debug, "debug", "d", false, "Enable debug mode")
	json2sqlCmd.Flags().StringVarP(&dbHost, "host", "H", "localhost", "DB host")
	json2sqlCmd.Flags().StringVarP(&dbPort, "port", "p", "5432", "DB port")
	json2sqlCmd.Flags().StringVarP(&dbUser, "username", "u", "admin", "DB user")
	json2sqlCmd.Flags().StringVarP(&dbPass, "password", "P", "admin", "DB password")
	json2sqlCmd.Flags().StringVarP(&dbSchema, "schema", "s", "gorawler", "DB schema")
	json2sqlCmd.Flags().StringVarP(&filePath, "model", "m", "", "Model json file path")
}

var (
	json2sqlCmd = &cobra.Command{
		Use:     "json2sql",
		Aliases: []string{"m"},
		Short:   "Import data to your SQL database from a nested json file",
		Run:     json2sql,
	}

	filePath string
	dbHost   string
	dbPort   string
	dbUser   string
	dbPass   string
	dbSchema string
)

func json2sql(cmd *cobra.Command, args []string) {
	var (
		err error
	)

	db, err := repository.NewDBWithConf(&repository.Config{
		Host:   dbHost,
		Port:   dbPort,
		User:   dbUser,
		Pass:   dbPass,
		DBName: dbSchema,
		Debug:  debug,
	})
	if err != nil {
		logrus.Errorln(err)
		return
	}

	c := core.New(db)

	err = c.Json2Sql(filePath)
	if err != nil {
		logrus.Errorln(err)
		return
	}
}
