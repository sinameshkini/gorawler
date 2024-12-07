package main

import (
	"github.com/sinameshkini/gorawler/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := cmd.Execute(); err != nil {
		logrus.Fatalln(err.Error())
	}
}
