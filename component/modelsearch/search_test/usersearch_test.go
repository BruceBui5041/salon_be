package modelsearch_test

// import (
// 	"bytes"
// 	"context"
// 	"log"
// 	"salon_be/component/logger"
// 	"salon_be/component/modelsearch"
// 	models "salon_be/model"
// 	"testing"

// 	"github.com/joho/godotenv"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// 	"gorm.io/gorm"
// )

// func TestDynamicFilterForUser(t *testing.T) {
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
// 		validateResult func(*testing.T, []models.User)
// 	}{
// 		{
// 			name:   "Filter services by price (greater than)",
// 			fields: []string{"created_services"},
// 			conditions: []map[string]interface{}{
// 				{"source": "created_services.id", "operator": "is not null", "target": ""},
// 			},
// 			validateResult: func(t *testing.T, services []models.User) {
// 				assert.NotEmpty(t, services)
// 				for _, c := range services {
// 					assert.NotEmpty(t, c.CreatedServices, uint64(50))
// 				}
// 			},
// 		},
// 	}

// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			var users []models.User

// 			var logBuffer bytes.Buffer
// 			log.SetOutput(&logBuffer)

// 			query := modelsearch.Search(
// 				ctx,
// 				appCtx.GetMainDBConnection().Model(&models.User{}),
// 				models.User{}.TableName(),
// 				tc.conditions,
// 				tc.fields,
// 				tc.orderBy,
// 			)

// 			sql := query.ToSQL(func(tx *gorm.DB) *gorm.DB {
// 				return tx.Find(&users)
// 			})
// 			t.Logf("Generated SQL: %s", sql)

// 			err := query.Find(&users).Error
// 			if err != nil {
// 				t.Logf("Error: %v", err)
// 			}

// 			t.Logf("Function logs:\n%s", logBuffer.String())

// 			require.NoError(t, err)
// 			tc.validateResult(t, users)
// 		})
// 	}
// }
