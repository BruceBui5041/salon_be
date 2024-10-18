package modelsearch_test

// import (
// 	"bytes"
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"
// 	"salon_be/component"
// 	"salon_be/component/logger"
// 	"salon_be/component/modelsearch"
// 	models "salon_be/model"
// 	"testing"
// 	"time"

// 	"github.com/joho/godotenv"
// 	"github.com/shopspring/decimal"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// 	gormlogger "gorm.io/gorm/logger"
// )

// func TestDynamicFilterForService(t *testing.T) {
// 	err := godotenv.Load("../../../.env")
// 	require.NoError(t, err, "Error loading .env file")
// 	ctx := context.Background()

// 	db, err := setupTestDB()
// 	require.NoError(t, err)

// 	appCtx := setupMockAppContext(db)

// 	logger.CreateAppLogger(ctx)

// 	var totalServices int64
// 	err = db.Model(&models.Service{}).Count(&totalServices).Error
// 	require.NoError(t, err)

// 	testCases := []struct {
// 		name           string
// 		fields         []string
// 		orderBy        *string
// 		conditions     interface{}
// 		validateResult func(*testing.T, []models.Service)
// 	}{
// 		{
// 			name:   "Filter services by price (greater than)",
// 			fields: []string{"price"},
// 			conditions: []map[string]interface{}{
// 				{"source": "price", "operator": ">", "target": 50},
// 			},
// 			validateResult: func(t *testing.T, services []models.Service) {
// 				assert.NotEmpty(t, services)
// 				for _, c := range services {
// 					assert.Greater(t, c.Price, uint64(50))
// 				}
// 			},
// 		},
// 		{
// 			name:       "get creator",
// 			fields:     []string{"creator"},
// 			conditions: []map[string]interface{}{},
// 			validateResult: func(t *testing.T, services []models.Service) {
// 				assert.NotEmpty(t, services)
// 				for _, c := range services {
// 					assert.NotEmpty(t, c.Creator.Email)
// 				}
// 			},
// 		},
// 		{
// 			name:   "Filter services by price > 50 OR creator_id = 6",
// 			fields: []string{"price", "creator_id"},
// 			conditions: []interface{}{
// 				map[string]interface{}{
// 					"source": "price", "operator": ">", "target": 50,
// 				},
// 				map[string]interface{}{
// 					"source": "creator_id", "operator": "=", "target": 6,
// 				},
// 			},
// 			validateResult: func(t *testing.T, services []models.Service) {
// 				assert.NotEmpty(t, services)
// 				for _, c := range services {
// 					assert.True(t, c.Price.GreaterThan(decimal.NewFromInt(50)) || c.CreatorID == 6,
// 						"Service should have price > 50 or creator_id = 6, but got price: %d and creator_id: %d",
// 						c.Price, c.CreatorID)
// 				}
// 			},
// 		},
// 		{
// 			name:   "Filter services by creator's email",
// 			fields: []string{"creator.email"},
// 			conditions: []interface{}{
// 				map[string]interface{}{
// 					"source": "creator.email", "operator": "=", "target": "abcd123@gmail.com",
// 				},
// 			},
// 			validateResult: func(t *testing.T, services []models.Service) {
// 				assert.NotEmpty(t, services)
// 				for _, c := range services {
// 					assert.Equal(t, "abcd123@gmail.com", c.Creator.Email,
// 						"Service creator's email should be abcd123@gmail.com, but got %s", c.Creator.Email)
// 				}
// 			},
// 		},
// 		{
// 			name:   "Filter services by creator's email AND price is 120",
// 			fields: []string{"creator.email"},
// 			conditions: []interface{}{
// 				map[string]interface{}{
// 					"source": "creator.email", "operator": "=", "target": "abcd123@gmail.com",
// 				},
// 				map[string]interface{}{
// 					"source": "price", "operator": "=", "target": 120,
// 				},
// 			},
// 			validateResult: func(t *testing.T, services []models.Service) {
// 				assert.NotEmpty(t, services)
// 				for _, c := range services {
// 					assert.Equal(t, "abcd123@gmail.com", c.Creator.Email,
// 						"Service creator's email should be abcd123@gmail.com, but got %s", c.Creator.Email,
// 					)
// 				}
// 			},
// 		},
// 		{
// 			name:   "Filter services by creator's email OR price is 120",
// 			fields: []string{"creator.email"},
// 			conditions: []interface{}{
// 				[]interface{}{
// 					map[string]interface{}{
// 						"source": "creator.email", "operator": "=", "target": "abcd123@gmail.com",
// 					},
// 				},

// 				[]interface{}{
// 					map[string]interface{}{
// 						"source": "id", "operator": "=", "target": 4,
// 					},
// 					map[string]interface{}{
// 						"source": "price", "operator": "=", "target": 120,
// 					},
// 				},
// 			},
// 			validateResult: func(t *testing.T, services []models.Service) {
// 				assert.NotEmpty(t, services)
// 				for _, c := range services {
// 					assert.Equal(t, "abcd123@gmail.com", c.Creator.Email,
// 						"Service creator's email should be abcd123@gmail.com, but got %s", c.Creator.Email,
// 					)

// 				}
// 			},
// 		},

// 		{
// 			name:   "creator.email == abcd123@gmail.com AND service.id = 4 AND price = 120 AND creator.role.name = instructor",
// 			fields: []string{"creator.email", "creator.roles.name", "price"},
// 			conditions: []interface{}{
// 				[]interface{}{
// 					map[string]interface{}{
// 						"source": "creator.email", "operator": "=", "target": "abcd123@gmail.com",
// 					},
// 					map[string]interface{}{
// 						"source": "creator.roles.name", "operator": "=", "target": "instructor",
// 					},
// 					map[string]interface{}{
// 						"source": "id", "operator": "=", "target": 4,
// 					},
// 					map[string]interface{}{
// 						"source": "price", "operator": "=", "target": 120,
// 					},
// 				},
// 			},
// 			validateResult: func(t *testing.T, services []models.Service) {
// 				assert.NotEmpty(t, services)
// 				for _, c := range services {
// 					assert.Equal(t, "abcd123@gmail.com", c.Creator.Email,
// 						"Service creator's email should be abcd123@gmail.com, but got %s", c.Creator.Email,
// 					)

// 					assert.Equal(t, uint64(120), c.Price,
// 						"Price should be 120, but got %d", c.Price,
// 					)

// 					assert.Equal(t, "instructor", c.Creator.Roles[0].Name,
// 						"Price should be 120, but got %d", c.Price,
// 					)
// 				}
// 			},
// 		},

// 		{
// 			name:   "service price between 80 and 125",
// 			fields: []string{"id", "price"},
// 			conditions: []interface{}{
// 				[]interface{}{
// 					map[string]interface{}{
// 						"source": "price", "operator": "between", "target": []interface{}{80, 125},
// 					},
// 					map[string]interface{}{
// 						"source": "id", "operator": "=", "target": 4,
// 					},
// 				},
// 			},
// 			validateResult: func(t *testing.T, services []models.Service) {
// 				assert.NotEmpty(t, services)
// 				for _, c := range services {
// 					assert.GreaterOrEqual(t, uint64(125), c.Price,
// 						"Price should be GreaterOrEqual 120, but got %d", c.Price,
// 					)
// 					assert.LessOrEqual(t, uint64(80), c.Price,
// 						"Price should be LessOrEqual 120, but got %d", c.Price,
// 					)
// 				}
// 			},
// 		},
// 		{
// 			name:   "service overview is null",
// 			fields: []string{"overview"},
// 			conditions: []interface{}{
// 				map[string]interface{}{
// 					"source": "overview", "operator": "is null", "target": nil,
// 				},
// 			},
// 			validateResult: func(t *testing.T, services []models.Service) {
// 				assert.NotEmpty(t, services)
// 				for _, c := range services {
// 					assert.Equal(t, "", c.Overview,
// 						"Overview should be NULL, but got %d", c.Overview,
// 					)
// 				}
// 			},
// 		},

// 		{
// 			name:   "service overview is not null",
// 			fields: []string{"overview"},
// 			conditions: []interface{}{
// 				map[string]interface{}{
// 					"source": "overview", "operator": "is not null", "target": nil,
// 				},
// 			},
// 			validateResult: func(t *testing.T, services []models.Service) {
// 				assert.NotEmpty(t, services)
// 				for _, c := range services {
// 					assert.NotEmpty(t, c.Overview,
// 						"Overview should be not NULL, but got %s", c.Overview,
// 					)
// 				}
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			var services []models.Service

// 			var logBuffer bytes.Buffer
// 			log.SetOutput(&logBuffer)

// 			query := modelsearch.Search(
// 				ctx,
// 				appCtx.GetMainDBConnection().Model(&models.Service{}),
// 				models.Service{}.TableName(),
// 				tc.conditions,
// 				tc.fields,
// 				tc.orderBy,
// 			)

// 			sql := query.ToSQL(func(tx *gorm.DB) *gorm.DB {
// 				return tx.Find(&services)
// 			})
// 			t.Logf("Generated SQL: %s", sql)

// 			err := query.Find(&services).Error
// 			if err != nil {
// 				t.Logf("Error: %v", err)
// 			}

// 			t.Logf("Function logs:\n%s", logBuffer.String())

// 			require.NoError(t, err)
// 			tc.validateResult(t, services)
// 		})
// 	}
// }

// func setupTestDB() (*gorm.DB, error) {
// 	dbUser := os.Getenv("DB_USER")
// 	dbPassword := os.Getenv("DB_PASSWORD")
// 	dbHost := os.Getenv("DB_HOST")
// 	dbPort := os.Getenv("DB_PORT")
// 	dbName := os.Getenv("DB_NAME")

// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
// 		dbUser, dbPassword, dbHost, dbPort, dbName)

// 	newLogger := gormlogger.New(
// 		log.New(os.Stdout, "\r\n", log.LstdFlags),
// 		gormlogger.Config{
// 			SlowThreshold:             time.Second,
// 			LogLevel:                  gormlogger.Silent,
// 			IgnoreRecordNotFoundError: true,
// 			Colorful:                  false,
// 		},
// 	)

// 	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
// 		Logger: newLogger,
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	return db, nil
// }

// func setupMockAppContext(db *gorm.DB) component.AppContext {
// 	return component.NewAppContext(
// 		db,
// 		nil, // Replace with mock implementations if needed
// 		nil,
// 		"test_jwt_secret",
// 		nil,
// 		nil,
// 		nil,
// 		nil,
// 		nil,
// 	)
// }
