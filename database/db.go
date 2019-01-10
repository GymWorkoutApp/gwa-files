package database

import (
	"context"
	"fmt"
	"github.com/jinzhu/gorm"
	"go.elastic.co/apm/module/apmgorm"
	_ "go.elastic.co/apm/module/apmgorm/dialects/postgres"
	"os"
	"sync"
	"sync/atomic"
)


var initialized uint32
var mu sync.Mutex
var instance *gorm.DB

func GetDB(context context.Context) *gorm.DB {
	if atomic.LoadUint32(&initialized) == 1 {
		instance = apmgorm.WithContext(context, instance)
		return instance
	}

	mu.Lock()
	defer mu.Unlock()

	if initialized == 0 {
		config := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_SSLMODE"))
		db, err := apmgorm.Open("postgres", config)
		if err != nil {
			fmt.Println(config)
			panic(err)
		}

		instance = db

		atomic.StoreUint32(&initialized, 1)
	}

	instance = apmgorm.WithContext(context, instance)

	return instance
}



