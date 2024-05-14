package main

import (
	"bash_server_project"
	"bash_server_project/pkg/handler"
	"bash_server_project/pkg/repository"
	"context"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

func InitConfigTest() error {
	viper.AddConfigPath("../configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func TestAppStartAndShutdown(t *testing.T) {
	// Тест инициализации конфигурации
	err := InitConfigTest()
	assert.NoError(t, err, "InitConfig should not return an error")

	// Тест загрузки переменных окружения
	err = godotenv.Load()
	assert.NoError(t, err, "Loading env variables should not return an error")

	// Тест инициализации базы данных
	repository.InitDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	assert.NotNil(t, repository.Db, "Database connection should be initialized")

	// Тест запуска HTTP-сервера
	srv := new(bash_server_project.Server)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = srv.Run(viper.GetString("port"), handler.InitRoutes())
		assert.NoError(t, err, "Running HTTP server should not return an error")
	}()

	// Ожидание запуска сервера
	time.Sleep(1 * time.Second)

	// Тест остановки HTTP-сервера
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	quit <- syscall.SIGTERM

	// Проверка остановки HTTP-сервера
	if srv != nil {
		err = srv.Shutdown(context.Background())
		assert.NoError(t, err, "Shutting down HTTP server should not return an error")
	}

	// Проверка закрытия соединения с базой данных
	if repository.Db != nil {
		err = repository.Db.Close()
		assert.NoError(t, err, "Closing database connection should not return an error")
	}
}
