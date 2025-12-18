package monoctl

import (
  "context"
  "errors"
  "fmt"
  "time"

  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

const mongoTimeout time.Duration = 10 * time.Second

type MongoDB struct {
  Uri string
  Collection string
  Fields []string
  QueryJSON string
  Limit int64

  ctx context.Context
}

func (m *MongoDB) validate() error {
  switch {
    case m.Uri == "":
      return errors.New("mongodb uri is required")
    case m.Collection == "":
      return errors.New("mongodb collection is required")
  }

  return nil
}

func (m *MongoDB) connect() (*mongo.Client, error) {
  ctxbg := context.Background()
  ctx, cancel := context.WithTimeout(ctxbg, mongoTimeout)
  defer cancel()

  client, err := mongo.Connect(ctx, options.Client().ApplyURI(m.Uri))
  if err != nil {
    return nil, fmt.Errorf("unable to connect: %w", err)
  }
  // defer client.Disconnect(ctx)

  if err := client.Ping(ctx, nil); err != nil {
    _ = client.Disconnect(ctx)
    return nil, fmt.Errorf("unable to ping: %w", err)
  }

  m.ctx = ctxbg
  return client, nil
}

func (m *MongoDB) Find(out any) error {
  if err := m.validate(); err != nil { return err }

  client, err := m.connect()
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

  if m.Limit > 0 {
    opts.SetLimit(m.Limit)
  }

  opts.SetSort(bson.M{"dateCreated": -1})

  // parse query
  var query bson.M
  if err := bson.UnmarshalExtJSON([]byte(m.QueryJSON), true, &query); err != nil {
    return fmt.Errorf("invalid query JSON: %w", err)
  }

  cs, err := connstring.ParseAndValidate(m.Uri)
  if err != nil {
    return fmt.Errorf("invalid mongodb uri: %w", err)
  }

  cur, err := client.Database(cs.Database).Collection(m.Collection).Find(m.ctx, query, opts)
  if err != nil {
    return fmt.Errorf("unable to run query: %w", err)
  }
  defer cur.Close(m.ctx)

  if err := cur.All(m.ctx, out); err != nil {
    return fmt.Errorf("unable to decode cursor: %w", err)
  }

  return nil
}

func (m *MongoDB) GetRows() ([][]any, error) {
  var rawData []map[string]any

  if err := m.Find(&rawData); err != nil {
    return nil, err
  }

  return mapToRows(m.Fields, rawData), nil
}
