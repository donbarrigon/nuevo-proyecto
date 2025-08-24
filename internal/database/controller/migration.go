package controller

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/migration"
)

func Migrare(ctx *app.HttpContext) {
	if !app.Env.DB_MIGRATION_ENABLE {
		ctx.ResponseError(app.Errors.Forbiddenf("Migration disabled"))
		app.PrintError("Migration disabled")
		panic("Migration disabled")
	}

	migration.Migrations = []app.List{}
	migration.Run()

	if er := os.MkdirAll(app.Env.LOG_PATH, os.ModePerm); er != nil {
		app.PrintError("Fail to create log directory :path: :error", app.E("path", app.Env.LOG_PATH), app.E("error", er.Error()))
		panic(er.Error())
	}

	filePath := filepath.Join(app.Env.LOG_PATH, "migration_tracker.txt")

	file, er := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if er != nil {
		app.PrintError("Fail to open file: :file :error", app.E("file", filePath), app.E("error", er.Error()))
		panic(er.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	records := []map[string]string{}

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		record := map[string]string{}
		for _, field := range fields {
			parts := strings.SplitN(field, ":", 2)
			if len(parts) < 2 {
				continue
			}
			record[parts[0]] = parts[1]
		}
		records = append(records, record)
	}
	if er := scanner.Err(); er != nil {
		app.PrintError("Fail to read file: :file :error", app.E("file", filePath), app.E("error", er.Error()))
		panic(er.Error())
	}

	migrations := []app.List{}
	for _, migration := range migration.Migrations {
		exists := false
		for _, record := range records {
			name := migration.Get("name").(string)
			if record["name"] == name {
				if record["action"] == "up" {
					exists = true
				}
				if record["action"] == "down" {
					exists = false
				}
			}
		}
		if !exists {
			migrations = append(migrations, migration)
		}
	}

	executedAt := time.Now()
	for _, migration := range migrations {
		migration.Get("up").(func())()
		line := fmt.Sprintf("executed_at:%s\taction:up\tname:%v\n", executedAt, migration.Get("name"))
		if _, er := file.WriteString(line); er != nil {
			app.PrintError("Fail to write :file :error", app.E("file", filePath), app.E("error", er.Error()))
			panic(er.Error())
		}
	}

	ctx.ResponseNoContent()
}

func Rollback(ctx *app.HttpContext) {
	if !app.Env.DB_MIGRATION_ENABLE {
		ctx.ResponseError(app.Errors.Forbiddenf("Migration disabled"))
		app.PrintError("Migration disabled")
		panic("Migration disabled")
	}

	migration.Migrations = []app.List{}
	migration.Run()

	if er := os.MkdirAll(app.Env.LOG_PATH, os.ModePerm); er != nil {
		app.PrintError("Fail to create log directory :path: :error", app.E("path", app.Env.LOG_PATH), app.E("error", er.Error()))
		panic(er.Error())
	}

	filePath := filepath.Join(app.Env.LOG_PATH, "migration_tracker.txt")

	file, er := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if er != nil {
		app.PrintError("Fail to open file: :file :error", app.E("file", filePath), app.E("error", er.Error()))
		panic(er.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	records := []map[string]string{}

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		record := map[string]string{}
		for _, field := range fields {
			parts := strings.SplitN(field, ":", 2)
			if len(parts) < 2 {
				continue
			}
			record[parts[0]] = parts[1]
		}
		records = append(records, record)
	}
	if er := scanner.Err(); er != nil {
		app.PrintError("Fail to read file: :file :error", app.E("file", filePath), app.E("error", er.Error()))
		panic(er.Error())
	}

	var recordsUp []map[string]string
	for _, record := range records {
		if record["action"] == "up" {
			recordsUp = append(recordsUp, record)
		}
		if record["action"] == "down" {
			for i, recordUp := range recordsUp {
				if recordUp["name"] == record["name"] {
					recordsUp = append(recordsUp[:i], recordsUp[i+1:]...)
				}
			}
		}
	}

	var last string
	if len(recordsUp) > 0 {
		last = recordsUp[len(recordsUp)-1]["executed_at"]
	}

	filtered := []map[string]string{}
	for _, record := range records {
		if record["executed_at"] == last {
			filtered = append(filtered, record)
		}
	}

	migrations := []app.List{}
	for _, filter := range filtered {
		for _, migration := range migration.Migrations {
			if migration.Get("name").(string) == filter["name"] {
				migrations = append(migrations, migration)
			}
		}
	}

	executedAt := time.Now()
	for _, migration := range migrations {
		migration.Get("down").(func())()
		line := fmt.Sprintf("executed_at:%s\taction:down\tname:%v\n", executedAt, migration.Get("name"))
		if _, er := file.WriteString(line); er != nil {
			app.PrintError("Fail to write :file :error", app.E("file", filePath), app.E("error", er.Error()))
			panic(er.Error())
		}
	}

	ctx.ResponseNoContent()
}
