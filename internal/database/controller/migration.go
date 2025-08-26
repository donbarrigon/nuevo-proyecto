package controller

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/migration"
)

func Migrare(ctx *app.HttpContext) {

	file := openFile("migration_tracker.txt")
	defer file.Close()

	migration.Migrations = []app.List{}
	migration.Run()

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
		app.PrintError("Fail to read file: :file :error", app.E("file", file.Name()), app.E("error", er.Error()))
		panic(er.Error())
	}

	migrations := []app.List{}
	for _, m := range migration.Migrations {
		exists := false
		for _, record := range records {
			name := m.Get("name").(string)
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
			migrations = append(migrations, m)
		}
	}
	runMigrations("up", migrations, file)

	app.PrintInfo("Migrations executed")
	ctx.ResponseNoContent()
}

func Rollback(ctx *app.HttpContext) {

	file := openFile("migration_tracker.txt")
	defer file.Close()

	migration.Migrations = []app.List{}
	migration.Run()

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
		if record["action"] == "up" {
			records = append(records, record)
		}
		if record["action"] == "down" {
			for i, recordUp := range records {
				if recordUp["name"] == record["name"] {
					records = append(records[:i], records[i+1:]...)
					break
				}
			}
		}
	}
	if er := scanner.Err(); er != nil {
		app.PrintError("Fail to read file: :file :error", app.E("file", file.Name()), app.E("error", er.Error()))
		panic(er.Error())
	}

	var last string
	if len(records) > 0 {
		last = records[len(records)-1]["executed_at"]
	}

	filtered := []map[string]string{}
	for _, record := range records {
		if record["executed_at"] == last {
			filtered = append(filtered, record)
		}
	}

	migrations := []app.List{}
	for _, filter := range filtered {
		for _, m := range migration.Migrations {
			if m.Get("name").(string) == filter["name"] {
				migrations = append(migrations, m)
			}
		}
	}
	runMigrations("down", migrations, file)

	app.PrintInfo("Migrations rolled back")
	ctx.ResponseNoContent()
}

func Reset(ctx *app.HttpContext) {

	file := openFile("migration_tracker.txt")
	defer file.Close()

	filePath := filepath.Join(app.Env.LOG_PATH, "seed_tracker.txt")
	er := os.Remove(filePath)
	if er != nil {
		if !os.IsNotExist(er) {
			fmt.Println("Fail to remove file seed_tracker:", er)
			return
		}
	}

	migration.Migrations = []app.List{}
	migration.Run()

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
		if record["action"] == "up" {
			records = append(records, record)
		}
		if record["action"] == "down" {
			for i, recordUp := range records {
				if recordUp["name"] == record["name"] {
					records = append(records[:i], records[i+1:]...)
					break
				}
			}
		}
	}
	if er := scanner.Err(); er != nil {
		app.PrintError("Fail to read file: :file :error", app.E("file", file.Name()), app.E("error", er.Error()))
		panic(er.Error())
	}

	migrations := []app.List{}
	for _, filter := range records {
		for _, m := range migration.Migrations {
			if m.Get("name").(string) == filter["name"] {
				migrations = append(migrations, m)
			}
		}
	}
	runMigrations("down", migrations, file)

	app.PrintInfo("Migration reset")
	ctx.ResponseNoContent()

}

func Refresh(ctx *app.HttpContext) {
	file := openFile("migration_tracker.txt")
	defer file.Close()

	filePath := filepath.Join(app.Env.LOG_PATH, "seed_tracker.txt")
	er := os.Remove(filePath)
	if er != nil {
		if !os.IsNotExist(er) {
			fmt.Println("Fail to remove file seed_tracker:", er)
			return
		}
	}

	migration.Migrations = []app.List{}
	migration.Run()

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
		if record["action"] == "up" {
			records = append(records, record)
		}
		if record["action"] == "down" {
			for i, recordUp := range records {
				if recordUp["name"] == record["name"] {
					records = append(records[:i], records[i+1:]...)
					break
				}
			}
		}
	}
	if er := scanner.Err(); er != nil {
		app.PrintError("Fail to read file: :file :error", app.E("file", file.Name()), app.E("error", er.Error()))
		panic(er.Error())
	}

	migrations := []app.List{}
	for _, filter := range records {
		for _, m := range migration.Migrations {
			if m.Get("name").(string) == filter["name"] {
				migrations = append(migrations, m)
			}
		}
	}
	runMigrations("down", migrations, file)
	runMigrations("up", migration.Migrations, file)

	app.PrintInfo("Migrations refreshed")
	ctx.ResponseNoContent()
}

func Fresh(ctx *app.HttpContext) {
	app.DB.Drop(context.TODO())

	filePath := filepath.Join(app.Env.LOG_PATH, "migration_tracker.txt")
	er := os.Remove(filePath)
	if er != nil {
		if !os.IsNotExist(er) {
			fmt.Println("Fail to remove file migration_tracker:", er)
			return
		}
	}

	filePath = filepath.Join(app.Env.LOG_PATH, "seed_tracker.txt")
	er = os.Remove(filePath)
	if er != nil {
		if !os.IsNotExist(er) {
			fmt.Println("Fail to remove file seed_tracker:", er)
			return
		}
	}

	file := openFile("migration_tracker.txt")
	defer file.Close()

	migration.Migrations = []app.List{}
	migration.Run()
	runMigrations("up", migration.Migrations, file)

	app.PrintInfo("Database refreshed")
	ctx.ResponseNoContent()
}

func openFile(fileName string) *os.File {
	if er := os.MkdirAll(app.Env.LOG_PATH, os.ModePerm); er != nil {
		app.PrintError("Fail to create log directory :path: :error", app.E("path", app.Env.LOG_PATH), app.E("error", er.Error()))
		panic(er.Error())
	}

	filePath := filepath.Join(app.Env.LOG_PATH, fileName)

	file, er := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if er != nil {
		app.PrintError("Fail to open file: :file :error", app.E("file", filePath), app.E("error", er.Error()))
		panic(er.Error())
	}
	return file
}

func runMigrations(action string, migrations []app.List, file *os.File) {
	executedAt := time.Now()
	for _, m := range migrations {
		m.Get(action).(func())()
		line := fmt.Sprintf("executed_at:%s\taction:%s\tname:%v\n", executedAt, action, m.Get("name"))
		if _, er := file.WriteString(line); er != nil {
			app.PrintError("Fail to write :file :error", app.E("file", file.Name()), app.E("error", er.Error()))
			panic(er.Error())
		}
	}
}
