package controller

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/donbarrigon/nuevo-proyecto/internal/app"
	"github.com/donbarrigon/nuevo-proyecto/internal/domain/db/seed"
)

func Seed(ctx app.HttpContext) {

	if !app.Env.SERVER_MIGRATION_ENABLE {
		ctx.ResponseError(app.Errors.Forbiddenf("Migration disabled"))
		return
	}

	seed.Run()

	if err := os.MkdirAll(app.Env.LOG_PATH, os.ModePerm); err != nil {
		ctx.ResponseError(app.Errors.InternalServerErrorf("fail to create log directory " + err.Error()))
		return
	}

	filePath := filepath.Join(app.Env.LOG_PATH, "seed_tracker.txt")

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		ctx.ResponseError(app.Errors.InternalServerErrorf("fail to open seed_tracker.txt file " + err.Error()))
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var records map[string]string

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

	if err := scanner.Err(); err != nil {
		ctx.ResponseError(app.Errors.InternalServerErrorf("Error leyendo el archivo: " + err.Error()))
		return
	}
	seeds := app.List{}
	for _, s := range seed.Seeds {
		if records[s.Key] == "" {
			seeds = append(seeds, s)
		}
	}

	for _, s := range seeds {
		s.Value.(func())()

		executedAt := time.Now().UTC().Format(time.RFC3339)
		line := fmt.Sprintf("executed_at:%s\tname:%s\n", executedAt, s.Key)
		if _, err := file.WriteString(line); err != nil {
			ctx.ResponseError(app.Errors.InternalServerErrorf("fail to write seed_tracker.txt " + err.Error()))
			return
		}
	}

	ctx.ResponseNoContent()

}
