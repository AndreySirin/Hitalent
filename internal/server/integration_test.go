package server

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	pgcontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"hitalent/internal/logger"
	"hitalent/internal/model"
	"hitalent/internal/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	username = "test"
	database = "test"
	passward = "test"
)

type Resp struct {
	Message string `json:"message"`
	ID      int    `json:"id"`
}

type IntegritionTest struct {
	suite.Suite
	container *pgcontainer.PostgresContainer
	handler   *Server
	id        uuid.UUID
}

func (suite *IntegritionTest) TestQuestCreateGood() {
	question := model.Question{
		Text: "тесты вообще работают?",
	}
	boby, err := json.Marshal(question)
	suite.NoError(err)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/questions/", bytes.NewBuffer(boby))
	request.Header.Set("Content-Type", "application/json")
	suite.handler.srv.Handler.ServeHTTP(recorder, request)

	var r Resp

	err = json.Unmarshal(recorder.Body.Bytes(), &r)
	suite.NoError(err)
	suite.NotEmpty(r.ID)
	suite.Equal("Вопрос успешно создан", r.Message)
	suite.Equal(http.StatusOK, recorder.Code)
}

func (suite *IntegritionTest) TestQuestCreateBad() {
	question := model.Question{
		Text: "?",
	}
	boby, err := json.Marshal(question)
	suite.NoError(err)

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/questions/", bytes.NewBuffer(boby))
	request.Header.Set("Content-Type", "application/json")
	suite.handler.srv.Handler.ServeHTTP(recorder, request)

	suite.Contains(recorder.Body.String(), "Text")
	suite.Equal(http.StatusBadRequest, recorder.Code)
}

func (suite *IntegritionTest) SetupSuite() {
	db, conteiner, err := InitPostgresBD()
	suite.Require().NoError(err)

	err = storage.MigrateUP(db)
	suite.Require().NoError(err)

	suite.container = conteiner
	suite.handler = New(logger.New(), validator.New(), "8080", db)

}

func (suite *IntegritionTest) TearDownSuite() {
	err := suite.container.Terminate(context.Background())
	suite.NoError(err)
}

func InitPostgresBD() (*gorm.DB, *pgcontainer.PostgresContainer, error) {
	ctx := context.Background()
	container, err := pgcontainer.Run(ctx, "postgres:17.4-alpine3.21",
		pgcontainer.WithDatabase(database),
		pgcontainer.WithUsername(username),
		pgcontainer.WithPassword(passward),
		pgcontainer.BasicWaitStrategies())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize postgres container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get container host: %w", err)
	}

	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get mapped port: %w", err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, mappedPort.Int(), username, passward, database,
	)

	stor, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize postgres storage: %w", err)
	}

	return stor, container, nil
}
func TestRunSuite(t *testing.T) {
	suite.Run(t, new(IntegritionTest))
}
