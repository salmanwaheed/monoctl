package monoctl

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

const timeout = 10 * time.Second

type dbConfig struct {
	client *mongo.Client
	uri    string
	ctx    context.Context
}

func NewDB(uri string) *dbConfig {
	return &dbConfig{uri: uri}
}

func (c *dbConfig) Connect() error {
	ctxBg := context.Background()
	ctx, cancel := context.WithTimeout(ctxBg, timeout)
	defer cancel()

	clientOptions := options.Client().ApplyURI(c.uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("connect error: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("ping error: %w", err)
	}

	c.client = client
	c.ctx = ctxBg

	log.Println("database is connected!")
	return nil
}

func (c *dbConfig) Disconnect() {
	if c.client != nil {
		_ = c.client.Disconnect(c.ctx)
		log.Println("database is disconnected!")
	}
}

type collection struct {
	col *mongo.Collection
	ctx context.Context
}

func (c *dbConfig) Collection(name string) *collection {
	cs, _ := connstring.ParseAndValidate(c.uri)

	return &collection{
		col: c.client.Database(cs.Database).Collection(name),
		ctx: c.ctx,
	}
}

func (c *collection) Aggregate(pipeline mongo.Pipeline) (Documents, error) {
	cursor, err := c.col.Aggregate(c.ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregate error: %w", err)
	}
	defer cursor.Close(c.ctx)

	var results Documents
	if err := cursor.All(c.ctx, &results); err != nil {
		return nil, fmt.Errorf("cursor decode error: %w", err)
	}

	return results, nil
}

type Query struct {
	Fields []string
	Category string
	DateCreated time.Time
}

func (q *Query) projection() bson.D {
	proj := bson.D{{Key: "_id", Value: 1}}

	for _, field := range q.Fields {
		proj = append(proj, bson.E{Key: field, Value: 1})
	}

	return proj
}

func (q *Query) Pipeline() mongo.Pipeline {
	return mongo.Pipeline{
		{{Key: "$match", Value: bson.M{
			"category": q.Category,
			"dateCreated": bson.M{"$gte": q.DateCreated},
		}}},
		{{Key: "$group", Value: bson.M{
			"_id": bson.M{"name": "$name", "email": "$email", "telephone": "$telephone"},
			"doc": bson.M{"$first": "$$ROOT"},
		}}},
		{{Key: "$replaceRoot", Value: bson.M{"newRoot": "$doc"}}},
		{{Key: "$project", Value: q.projection()}},
		{{Key: "$sort", Value: bson.M{"dateCreated": -1}}},
		// {{Key: "$limit", Value: 2}},
	}
}

type Documents []map[string]any

func formatValue(v any) string {
	switch t := v.(type) {
		case nil:
			return ""
		case primitive.DateTime:
			return t.Time().Format("2006-01-02 15:04:05")
		default:
			return fmt.Sprintf("%v", v)
	}
}

// Convert "[]map" to "[][]any" for google sheets
func (docs Documents) ToRows(fields []string) [][]any {
	var rows [][]any

	for _, doc := range docs {
		row := make([]any, len(fields))
		for i, key := range fields {
			row[i] = formatValue(doc[key])
		}
		rows = append(rows, row)
	}

	return rows
}

// Convert "[]map" to "[][]string" for CSV
func (docs Documents) ToCSV(fields []string) [][]string {
	var rows [][]string

	for _, doc := range docs {
		row := make([]string, len(fields))
		for i, key := range fields {
			row[i] = formatValue(doc[key])
		}
		rows = append(rows, row)
	}

	return rows
}
