package main

import (
	"database/sql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
	"yellot-parser/main/instance"
	"yellot-parser/main/logger"
	"yellot-parser/main/models"
)

func main() {
	consoleLogger := &logger.ConsoleLogger{}

	sqlDb, err := sql.Open("pgx", "host=localhost user=postgres password=postgres port=5432 sslmode=disable") // dev

	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDb,
	}), &gorm.Config{})

	if err != nil {
		consoleLogger.Warning(err.Error())
		return
	}

	// dev
	_ = gormDb.Migrator().DropTable(&models.Resource{})
	_ = gormDb.Migrator().DropTable(&models.Order{})
	_ = gormDb.Migrator().DropTable(&models.Suggestion{})
	_ = gormDb.Migrator().DropTable(&models.ParseFile{})
	_ = gormDb.Migrator().DropTable(&models.Provider{})

	err = gormDb.AutoMigrate(&models.Order{})
	if err != nil {
		consoleLogger.Warning(err.Error())
		return
	}
	err = gormDb.AutoMigrate(&models.ParseFile{})
	if err != nil {
		consoleLogger.Warning(err.Error())
		return
	}
	err = gormDb.AutoMigrate(&models.Provider{})
	if err != nil {
		consoleLogger.Warning(err.Error())
		return
	}
	err = gormDb.AutoMigrate(&models.Resource{})
	if err != nil {
		consoleLogger.Warning(err.Error())
		return
	}
	err = gormDb.AutoMigrate(&models.Suggestion{})
	if err != nil {
		consoleLogger.Warning(err.Error())
		return
	}

	metalPortalParser := instance.MetalPortalParser(consoleLogger, gormDb)
	metaloobrabotchikiParser := instance.MetaloobrabotchikiParser(consoleLogger, gormDb)
	obrabotkanetParser := instance.ObrabotkaNetParser(consoleLogger, gormDb)
	obrabotkanetArchiveParser := instance.ObrabotkaNetArchiveParser(consoleLogger, gormDb)
	prommarketParser := instance.PromMarketParser(consoleLogger, gormDb)

	var wg sync.WaitGroup
	go metalPortalParser.Parse(&wg)
	wg.Add(1)
	go metaloobrabotchikiParser.Parse(&wg)
	wg.Add(1)
	go obrabotkanetParser.Parse(&wg)
	wg.Add(1)
	go obrabotkanetArchiveParser.Parse(&wg)
	wg.Add(1)
	go prommarketParser.Parse(&wg)
	wg.Add(1)
	wg.Wait()

	consoleLogger.Info("Parsing finished")
}
