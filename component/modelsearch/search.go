package modelsearch

import (
	"context"
	"fmt"
	"reflect"
	"salon_be/common"
	"salon_be/component/logger"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type FieldPreload struct {
	OrderBy string `json:"order_by"`
	Field   string `json:"field"`
}

func Search(
	ctx context.Context,
	db *gorm.DB,
	model string,
	condition interface{},
	fields interface{},
	orderBy *string, // ex: "id desc" or "id asc"
) *gorm.DB {
	db = applyJoins(ctx, db, condition, model)
	if orderBy != nil {
		db = applyOrderBy(db, *orderBy)
	}
	db = applyFieldPreloads(db, fields)
	// db = applyConditionPreloads(db, condition)

	switch reflect.TypeOf(condition).Kind() {
	case reflect.Slice:
		return applySliceFilter(ctx, db, condition, model)
	default:
		return db
	}
}

func applyFieldPreloads(db *gorm.DB, fields interface{}) *gorm.DB {
	preloaded := make(map[string]bool)

	modelType := reflect.TypeOf(db.Statement.Model)
	for modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	switch reflect.TypeOf(fields).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(fields)
		for i := 0; i < s.Len(); i++ {
			field := s.Index(i).Interface()
			switch f := field.(type) {
			case string:
				fieldPreload := FieldPreload{
					Field: f,
				}
				db = processStringField(db, fieldPreload, modelType, preloaded)
			case map[string]interface{}:
				fieldPreload := FieldPreload{}
				if fieldStr, ok := f["field"].(string); ok {
					fieldPreload.Field = fieldStr
				}

				if orderBy, ok := f["order_by"].(string); ok {
					fieldPreload.OrderBy = orderBy
				}

				db = processStringField(db, fieldPreload, modelType, preloaded)
			default:
				// handle other types if needed
			}
		}
	default:
		// Handle the case where fields is not a slice
	}

	return db
}

func processStringField(db *gorm.DB, fieldStr FieldPreload, modelType reflect.Type, preloaded map[string]bool) *gorm.DB {
	parts := strings.Split(fieldStr.Field, ".")
	orderBy := fieldStr.OrderBy
	var preloadPath []string
	currentType := modelType

	for i, part := range parts {
		fieldName := getStructFieldName(currentType, part)
		if fieldName != "" {
			field, found := currentType.FieldByName(fieldName)
			if !found {
				break
			}

			if isRelation(field.Tag.Get("gorm")) {
				preloadPath = append(preloadPath, fieldName)
				preloadString := strings.Join(preloadPath, ".")
				if !preloaded[preloadString] {
					db = db.Preload(preloadString, func(d *gorm.DB) *gorm.DB {
						// using the last struct field to order
						if i != len(parts)-1 {
							return d
						}

						if orderBy != "" {
							d = applyOrderBy(d, orderBy)
						}

						return d
					})
					preloaded[preloadString] = true
				}
				currentType = field.Type
				for currentType.Kind() == reflect.Ptr {
					currentType = currentType.Elem()
				}
				if currentType.Kind() == reflect.Slice {
					currentType = currentType.Elem()
					for currentType.Kind() == reflect.Ptr {
						currentType = currentType.Elem()
					}
				}
			} else {
				break
			}
		} else {
			break
		}
	}
	return db
}

func isRelation(gormTag string) bool {
	return strings.Contains(gormTag, "many2many") ||
		strings.Contains(gormTag, "has many") ||
		strings.Contains(gormTag, "has one") ||
		strings.Contains(gormTag, "belongs to") ||
		strings.Contains(gormTag, "foreignKey") ||
		strings.Contains(gormTag, "references") ||
		strings.Contains(gormTag, "constraint")
}

func applyJoins(ctx context.Context, db *gorm.DB, condition interface{}, model string) *gorm.DB {
	joinedTables := make(map[string]bool)

	var applyJoinsRecursive func(interface{}, reflect.Type, string, []string)
	applyJoinsRecursive = func(f interface{}, modelType reflect.Type, currentTable string, path []string) {
		s := reflect.ValueOf(f)
		if s.Kind() != reflect.Slice {
			return
		}

		for i := 0; i < s.Len(); i++ {
			element := s.Index(i).Interface()
			switch reflect.TypeOf(element).Kind() {
			case reflect.Map:
				m := element.(map[string]interface{})
				source := m["source"].(string)
				parts := strings.Split(source, ".")

				if len(parts) > 1 {
					var err error
					db, err = handleNestedRelationship(ctx, db, modelType, parts[:len(parts)-1], currentTable, joinedTables)
					if err != nil {
						db.AddError(err)
						return
					}
				}
			case reflect.Slice:
				applyJoinsRecursive(element, modelType, currentTable, path)
			}
		}
	}

	modelType := reflect.TypeOf(db.Statement.Model)
	for modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}
	applyJoinsRecursive(condition, modelType, model, []string{})
	return db
}

func handleNestedRelationship(
	ctx context.Context,
	db *gorm.DB,
	modelType reflect.Type,
	parts []string,
	currentTable string,
	joinedTables map[string]bool,
) (*gorm.DB, error) {
	for i, part := range parts {
		field, ok := findField(modelType, part)
		if !ok {
			return db, fmt.Errorf("field related to '%s' not found in %s", part, modelType.Name())
		}
		relationType := field.Type
		for relationType.Kind() == reflect.Ptr || relationType.Kind() == reflect.Slice {
			relationType = relationType.Elem()
		}
		relationTable := getTableName(relationType)

		joinKey := fmt.Sprintf("%s_%s", currentTable, relationTable)
		if !joinedTables[joinKey] {
			if field.Type.Kind() == reflect.Slice {
				joinTable, _, joinForeignKey, references, joinReferences := getManyToManyInfo(modelType, field)
				manyJoinClause := fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.%s LEFT JOIN %s ON %s.%s = %s.%s",
					joinTable, currentTable, references, joinTable, joinForeignKey,
					relationTable, joinTable, joinReferences, relationTable, references)
				db = db.Joins(manyJoinClause)
			} else {
				foreignKey := getForeignKeyName(field)
				joinClause := fmt.Sprintf("LEFT JOIN %s ON %s.%s = %s.id",
					relationTable, currentTable, foreignKey, relationTable)
				db = db.Joins(joinClause)
			}
			joinedTables[joinKey] = true
		}

		if i == len(parts)-1 {
			break
		}

		currentTable = relationTable
		modelType = relationType
	}
	return db, nil
}

func findField(t reflect.Type, name string) (reflect.StructField, bool) {
	if t.Kind() == reflect.Slice || t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return reflect.StructField{}, false
	}

	parts := strings.Split(name, ".")
	for i, part := range parts {
		found := false
		for j := 0; j < t.NumField(); j++ {
			field := t.Field(j)
			tag := field.Tag.Get("json")
			jsonName := strings.Split(tag, ",")[0]
			if jsonName == part {
				if i == len(parts)-1 {
					return field, true
				}
				t = field.Type
				if t.Kind() == reflect.Slice || t.Kind() == reflect.Ptr {
					t = t.Elem()
				}
				found = true
				break
			}
		}
		if !found {
			return reflect.StructField{}, false
		}
	}
	return reflect.StructField{}, false
}

func getManyToManyInfo(modelType reflect.Type, field reflect.StructField) (string, string, string, string, string) {
	gormTag := field.Tag.Get("gorm")
	joinTable := ""
	foreignKey := ""
	joinForeignKey := ""
	references := ""
	joinReferences := ""

	for _, tag := range strings.Split(gormTag, ";") {
		if strings.HasPrefix(tag, "many2many:") {
			parts := strings.Split(tag, ":")
			if len(parts) > 1 {
				joinTable = parts[1]
			}
		}

		if strings.HasPrefix(tag, "foreignKey:") {
			foreignKey = toSnakeCase(strings.TrimPrefix(tag, "foreignKey:"))
		} else if strings.HasPrefix(tag, "references:") {
			references = toSnakeCase(strings.TrimPrefix(tag, "references:"))
		} else if strings.HasPrefix(tag, "joinForeignKey:") {
			joinForeignKey = toSnakeCase(strings.TrimPrefix(tag, "joinForeignKey:"))
		} else if strings.HasPrefix(tag, "joinReferences:") {
			joinReferences = toSnakeCase(strings.TrimPrefix(tag, "joinReferences:"))
		}
	}

	if joinTable == "" {
		joinTable = fmt.Sprintf("%s_%s", toSnakeCase(modelType.Name()), toSnakeCase(field.Name))
	}
	if foreignKey == "" {
		foreignKey = toSnakeCase(modelType.Name()) + "_id"
	}
	if references == "" {
		references = "id"
	}
	if joinForeignKey == "" {
		joinForeignKey = foreignKey
	}
	if joinReferences == "" {
		joinReferences = references
	}

	return joinTable, foreignKey, joinForeignKey, references, joinReferences
}

func applySliceFilter(ctx context.Context, db *gorm.DB, condition interface{}, model string) *gorm.DB {
	s := reflect.ValueOf(condition)
	orConditions := make([]string, 0)
	orArgs := make([]interface{}, 0)

	for i := 0; i < s.Len(); i++ {
		element := s.Index(i).Interface()
		switch reflect.TypeOf(element).Kind() {
		case reflect.Slice:
			subConditions, subArgs := buildConditions(ctx, db, element, model)
			if len(subConditions) > 0 {
				orConditions = append(orConditions, "("+strings.Join(subConditions, " AND ")+")")
				orArgs = append(orArgs, subArgs...)
			}
		case reflect.Map:
			condition, args := buildCondition(ctx, db, element.(map[string]interface{}), model)
			if condition != "" {
				orConditions = append(orConditions, condition)
				orArgs = append(orArgs, args...)
			}
		}
	}

	if len(orConditions) > 0 {
		return db.Where(strings.Join(orConditions, " OR "), orArgs...)
	}
	return db
}

func buildConditions(ctx context.Context, db *gorm.DB, condition interface{}, model string) ([]string, []interface{}) {
	s := reflect.ValueOf(condition)
	conditions := make([]string, 0)
	args := make([]interface{}, 0)

	for i := 0; i < s.Len(); i++ {
		element := s.Index(i).Interface()
		switch reflect.TypeOf(element).Kind() {
		case reflect.Map:
			condition, arg := buildCondition(ctx, db, element.(map[string]interface{}), model)
			if condition != "" {
				conditions = append(conditions, condition)
				args = append(args, arg...)
			}
		case reflect.Slice:
			subConditions, subArgs := buildConditions(ctx, db, element, model)
			if len(subConditions) > 0 {
				conditions = append(conditions, "("+strings.Join(subConditions, " OR ")+")")
				args = append(args, subArgs...)
			}
		}
	}

	return conditions, args
}

func buildCondition(ctx context.Context, db *gorm.DB, m map[string]interface{}, model string) (string, []interface{}) {
	source := m["source"].(string)
	operator := m["operator"].(string)
	target := m["target"]

	column := getColumn(db, source, model)

	// Special handling for "id" source
	if strings.HasSuffix(strings.ToLower(source), "id") {
		return buildIDCondition(ctx, db, source, model, column, operator, target)
	}

	switch strings.ToLower(operator) {
	case "=", "!=", ">", "<", ">=", "<=":
		return fmt.Sprintf("%s %s ?", column, operator), []interface{}{target}
	case "like":
		return fmt.Sprintf("%s LIKE ?", column), []interface{}{fmt.Sprintf("%%%v%%", target)}
	case "in":
		return fmt.Sprintf("%s IN (?)", column), []interface{}{target}
	case "not in":
		return fmt.Sprintf("%s NOT IN (?)", column), []interface{}{target}
	case "between":
		targetSlice, ok := target.([]interface{})
		if !ok || len(targetSlice) != 2 {
			logger.AppLogger.Error(ctx,
				"between value must be a slice with 2 items",
				zap.Any("currentTarget", target),
			)
			return "", nil // Invalid input for BETWEEN
		}
		return fmt.Sprintf("%s BETWEEN ? AND ?", column), targetSlice
	case "is null":
		return fmt.Sprintf("%s IS NULL", column), []interface{}{}
	case "is not null":
		return fmt.Sprintf("%s IS NOT NULL", column), []interface{}{}
	default:
		logger.AppLogger.Warn(ctx, "Unsupported operator", zap.String("operator", operator))
		return "", nil
	}
}

func buildIDCondition(
	ctx context.Context,
	db *gorm.DB,
	source string,
	model string,
	column,
	operator string,
	target interface{},
) (string, []interface{}) {
	parseID := func(id string) uint32 {
		uid, err := common.FromBase58(id)
		if err != nil {
			panic(common.ErrInvalidRequest(err))
		}
		return uid.GetLocalID()
	}

	switch strings.ToLower(operator) {
	case "=", "!=", ">", "<", ">=", "<=":
		parsedID := parseID(target.(string))
		return fmt.Sprintf("%s %s ?", column, operator), []interface{}{parsedID}
	case "in", "not in":
		targets, ok := target.([]interface{})
		if !ok {
			logger.AppLogger.Error(ctx, "IN operator for id must use a slice")
			return "", nil
		}

		ids := make([]string, len(targets))
		for i, v := range targets {
			valueStr, ok := v.(string)
			if !ok {
				logger.AppLogger.Error(ctx, "IN operator for id must use a slice of strings")
				return "", nil
			}
			ids[i] = valueStr
		}

		parsedIDs := make([]uint32, len(ids))
		for i, id := range ids {
			parsedIDs[i] = parseID(id)
		}
		return fmt.Sprintf("%s IN (?)", column), []interface{}{parsedIDs}
	case "between":
		ids, ok := target.([]string)
		if !ok || len(ids) != 2 {
			logger.AppLogger.Error(ctx, "BETWEEN operator for id must use a slice of 2 strings")
			return "", nil
		}
		parsedIDs := []interface{}{parseID(ids[0]), parseID(ids[1])}
		return fmt.Sprintf("%s BETWEEN ? AND ?", column), parsedIDs
	default:
		return buildCondition(ctx, db, map[string]interface{}{
			"source":   source,
			"operator": operator,
			"target":   target,
		}, model)
	}
}

func getColumn(db *gorm.DB, source string, model string) string {
	parts := strings.Split(source, ".")
	if len(parts) > 1 {
		modelType := reflect.TypeOf(db.Statement.Model).Elem()
		currentTable := getTableName(modelType)

		for i, part := range parts[:len(parts)-1] {
			field, ok := findField(modelType, part)
			if ok {
				relationType := field.Type
				for relationType.Kind() == reflect.Ptr || relationType.Kind() == reflect.Slice {
					relationType = relationType.Elem()
				}

				currentTable = getTableName(relationType)
				modelType = relationType
			} else if i == 0 {
				currentTable = model
			} else {
				currentTable = toSnakeCase(part)
			}
		}

		lastField, ok := findField(modelType, parts[len(parts)-1])
		if ok {
			return fmt.Sprintf("%s.%s", currentTable, toSnakeCase(lastField.Name))
		}
		return fmt.Sprintf("%s.%s", currentTable, toSnakeCase(parts[len(parts)-1]))
	}

	return fmt.Sprintf("%s.%s", model, toSnakeCase(source))
}

func getStructFieldName(t reflect.Type, jsonField string) string {
	originalType := t

	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}

	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		if originalType.Kind() == reflect.Slice && originalType.Elem().Kind() == reflect.Ptr {
			return jsonField
		}
		return ""
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("json")
		name := strings.Split(tag, ",")[0]
		if name == jsonField {
			return field.Name
		}
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if strings.EqualFold(field.Name, jsonField) {
			return field.Name
		}
	}

	if field, found := t.FieldByName(jsonField); found {
		return field.Name
	}

	return ""
}
func getTableName(t reflect.Type) string {
	if method, ok := t.MethodByName("TableName"); ok {
		results := method.Func.Call([]reflect.Value{reflect.New(t).Elem()})
		if len(results) > 0 {
			return results[0].String()
		}
	}
	return toSnakeCase(t.Name())
}

func getForeignKeyName(field reflect.StructField) string {
	gormTag := field.Tag.Get("gorm")
	for _, tag := range strings.Split(gormTag, ";") {
		if strings.HasPrefix(tag, "foreignKey:") {
			return toSnakeCase(strings.TrimPrefix(tag, "foreignKey:"))
		}
	}
	return toSnakeCase(field.Name) + "_id"
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 {
			if r >= 'A' && r <= 'Z' {
				if (s[i-1] != 'I' || i < len(s)-1 && s[i+1] != 'D') &&
					(i < 2 || s[i-2:i] != "ID") {
					result.WriteRune('_')
				}
			}
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
