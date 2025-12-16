package monoctl

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type MongoDB struct {
  Uri string
  Collection string
  Fields []string
  QueryJSON string

  ctx context.Context
  timeout time.Duration
  // client *mongo.Client
}

func (m *MongoDB) auth() (*mongo.Client, error) {
  m.timeout = 10 * time.Second

  ctx, cancel := context.WithTimeout(context.TODO(), m.timeout)
  defer cancel()

  client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.Uri))
  if err != nil {
    return nil, fmt.Errorf("unable to connect: %v", err)
  }
  // defer client.Disconnect(ctx)

  if err := client.Ping(ctx, nil); err != nil {
    return nil, fmt.Errorf("unable to ping: %v", err)
  }

  m.ctx = context.TODO()
  return client, nil
}

func (m *MongoDB) Find(out any) error {
  client, err := m.auth()
  if err != nil {
    return err
  }
  defer client.Disconnect(m.ctx)

  // parse projection
  // projMap := bson.M{}
  // var fieldsJSON string = `["name","email"]`
  // var fields []string
  // if err := json.Unmarshal([]byte(fieldsJSON), &fields); err != nil {
  //   for _, f := range fields {
  //     projMap[f] = 1
  //   }
  // }

  opts := options.Find()
  // opts.SetProjection(projMap)
  opts.SetLimit(2)
  opts.SetSort(bson.M{"dateCreated": -1})

  // parse query
  var query bson.M
  bson.UnmarshalExtJSON([]byte(m.QueryJSON), true, &query)

  cs, _ := connstring.ParseAndValidate(m.Uri)
  cur, err := client.Database(cs.Database).Collection(m.Collection).Find(m.ctx, query, opts)
  if err != nil {
    return fmt.Errorf("unable to run query: %v", err)
  }
  defer cur.Close(m.ctx)

  // var rawData []map[string]any
  // if err := cur.All(m.ctx, &rawData); err != nil {
  //   return fmt.Errorf("unable to decode cursor: %v", err)
  // }
  // fmt.Println(rawData)

  return cur.All(m.ctx, out)
}

func (m *MongoDB) GetRows() ([][]any, error) {
  var rawData []map[string]any

  if err := m.Find(&rawData); err != nil {
    return nil, fmt.Errorf("unable to decode cursor: %v", err)
  }

  return mapToRows(m.Fields, rawData), nil
}
