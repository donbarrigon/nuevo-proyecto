package controller

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/database/seed"
)

func SeedRun(ctx *app.HttpContext) {

	if !app.Env.SERVER_MIGRATION_ENABLE {
		ctx.ResponseError(app.Errors.Forbiddenf("Migration disabled"))
		app.PrintError("Migration disabled")
		panic("Migration disabled")
	}

	seed.Seeds = app.List{}
	seed.Run()

	if er := os.MkdirAll(app.Env.LOG_PATH, os.ModePerm); er != nil {
		app.PrintError("Fail to create log directory :path: :error", app.E("path", app.Env.LOG_PATH), app.E("error", er.Error()))
		panic(er.Error())
	}

	filePath := filepath.Join(app.Env.LOG_PATH, "seed_tracker.txt")

	file, er := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if er != nil {
		app.PrintError("Fail to open file: :file :error", app.E("file", filePath), app.E("error", er.Error()))
		panic(er.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	records := map[string]string{}

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		var key string
		var value string
		for _, field := range fields {
			parts := strings.SplitN(field, ":", 2)
			if len(parts) == 2 {
				if parts[0] == "name" {
					key = parts[1]
				} else {
					value = parts[1]
				}
			}
		}
		records[key] = value
	}

	if er := scanner.Err(); er != nil {
		app.PrintError("Fail to read file: :file :error", app.E("file", filePath), app.E("error", er.Error()))
		panic(er.Error())
	}
	seeds := app.List{}
	for _, s := range seed.Seeds {
		if records[s.Key] == "" {
			seeds = append(seeds, s)
		}
	}

	results := map[string]string{}
	for _, s := range seeds {
		s.Value.(func())()
		executedAt := time.Now().UTC().Format(time.RFC3339)
		results[s.Key] = executedAt
		line := fmt.Sprintf("executed_at:%s\tname:%s\n", executedAt, s.Key)

		if _, er := file.WriteString(line); er != nil {
			app.PrintError("Fail to write :file :error", app.E("file", filePath), app.E("error", er.Error()))
			panic(er.Error())
		}
	}

	ctx.ResponseNoContent()

}

func SeedList(ctx *app.HttpContext) {
	if !app.Env.SERVER_MIGRATION_ENABLE {
		ctx.ResponseError(app.Errors.Forbiddenf("Migration disabled"))
		app.PrintError("Migration disabled")
		panic("Migration disabled")
	}

	if er := os.MkdirAll(app.Env.LOG_PATH, os.ModePerm); er != nil {
		app.PrintError("Fail to create log directory :path: :error", app.E("path", app.Env.LOG_PATH), app.E("error", er.Error()))
		panic(er.Error())
	}

	filePath := filepath.Join(app.Env.LOG_PATH, "seed_tracker.txt")

	file, er := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if er != nil {
		app.PrintError("Fail to open :file: :error", app.E("file", filePath), app.E("error", er.Error()))
		panic(er.Error())
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	records := map[string]string{}

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		var key string
		var value string
		for _, field := range fields {
			parts := strings.SplitN(field, ":", 2)
			if len(parts) == 2 {
				if parts[0] == "name" {
					key = parts[1]
				} else {
					value = parts[1]
				}
			}
		}
		records[key] = value
	}

	if er := scanner.Err(); er != nil {
		app.PrintError("Fail to read file: :file :error", app.E("file", filePath), app.E("error", er.Error()))
		panic(er.Error())
	}

	app.PrintInfo("Seed tracker :records", app.E("records", records))

	ctx.ResponseNoContent()
}

func SeedForce(ctx *app.HttpContext) {
	if !app.Env.SERVER_MIGRATION_ENABLE {
		ctx.ResponseError(app.Errors.Forbiddenf("Migration disabled"))
		panic("Migration disabled")
	}

	seed.Seeds = app.List{}
	seed.Run()

	if er := os.MkdirAll(app.Env.LOG_PATH, os.ModePerm); er != nil {
		app.PrintError("Fail to create log directory :path: :error", app.E("path", app.Env.LOG_PATH), app.E("error", er.Error()))
		panic(er.Error())
	}

	filePath := filepath.Join(app.Env.LOG_PATH, "seed_tracker.txt")

	file, er := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if er != nil {
		app.PrintError("Fail to open :file: :error", app.E("file", filePath), app.E("error", er.Error()))
		panic(er.Error())
	}
	defer file.Close()

	seed.Seeds.Get(ctx.Params["name"]).(func())()

	executedAt := time.Now().UTC().Format(time.RFC3339)
	line := fmt.Sprintf("executed_at:%s\tname:%s\n", executedAt, ctx.Params["name"])

	if _, er := file.WriteString(line); er != nil {
		app.PrintError("Fail to write :file :error", app.E("file", filePath), app.E("error", er.Error()))
		panic(er.Error())
	}
	ctx.ResponseNoContent()
}
