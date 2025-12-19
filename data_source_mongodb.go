package monoctl

import (
  "context"
  "encoding/json"
  "errors"
  "fmt"
  "time"

  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "go.mongodb.org/mongo-driver/mongo/readpref"
  "go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

const mongoTimeout time.Duration = 30 * time.Second

type MongoDB struct {
  Uri string
  Collection string
  Query QueryDSL

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

  opts := options.Client().ApplyURI(m.Uri).SetMinPoolSize(5).SetMaxPoolSize(100)
  client, err := mongo.Connect(ctx, opts)
  if err != nil {
    return nil, fmt.Errorf("unable to connect: %w", err)
  }
  // defer client.Disconnect(ctx)

  if err := client.Ping(ctx, readpref.Primary()); err != nil {
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

  if m.Query.Limit > 0 {
    opts.SetLimit(m.Query.Limit)
  }

  if len(m.Query.Sort) > 0 {
    sort, err := m.Query.Sort.ToBSON()
    if err != nil { return err }

    opts.SetSort(sort)
  }

  query, err := m.Query.Filter.ToBSON()
  if err != nil { return err }

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

  return mapToRows(m.Query.Select, rawData), nil
}

type sort map[string]string

func (s *sort) String() string {
  if s == nil || len(*s) == 0 {
    return ""
  }

  b, err := json.Marshal(*s)
  if err != nil {
    return fmt.Sprintf("failed to marshal sort: %v", err)
  }

  return string(b)
}

func (s *sort) Set(v string) error {
  if v == "" {
    return errors.New("sort input is empty")
  }

  if err := json.Unmarshal([]byte(v), s); err != nil {
    return fmt.Errorf("invalid sort json: %v", err)
  }

  return nil
}

func (s *sort) Type() string {
  return "sort"
}

func (s *sort) ToBSON() (bson.D, error) {
  out := bson.D{}

  for field, val := range *s {
    switch val {
      case "asc", "1":
        out = append(out, bson.E{Key: field, Value: 1})
      case "desc", "-1":
        out = append(out, bson.E{Key: field, Value: -1})
      default:
        return nil, fmt.Errorf("unsupported operator %q for field %q", val, field)
    }
  }

  return out, nil
}

type filter map[string]any

var allowedFilterOperators = map[string]struct{}{
  "eq":  {},
  "ne":  {},
  "gt":  {},
  "gte": {},
  "lt":  {},
  "lte": {},
  "in":  {},
  "nin": {},
  "regex": {},
  "exists": {},
}

func (f *filter) String() string {
  if f == nil || len(*f) == 0 {
    return ""
  }

  b, err := json.Marshal(*f)
  if err != nil {
    return fmt.Sprintf("failed to marshal filter: %v", err)
  }

  return string(b)
}

func (f *filter) Set(v string) error {
  if v == "" {
    return errors.New("filter input is empty")
  }

  if err := json.Unmarshal([]byte(v), f); err != nil {
    return fmt.Errorf("invalid filter json: %v", err)
  }

  return nil
}

func (f *filter) Type() string {
  return "filter"
}

func (f *filter) ToBSON() (bson.M, error) {
  out := bson.M{}

  for field, raw := range *f {
    switch v := raw.(type) {
      // case 1: implicit equality
      case string, int, int64, float64, bool:
        out[field] = v

      // case 2: operator based
      case filter:
        opMap := bson.M{}

        for op, val := range v {
          if _, ok := allowedFilterOperators[op]; !ok {
            return nil, fmt.Errorf("unsupported operator %q for field %q", op, field)
          }

          // convert "2025-12-19T00:00:00+04:00" as string to time.Time
          if s, ok := val.(string); ok {
            if t, err := time.Parse(time.RFC3339, s); err == nil {
              val = t
            }
          }

          opMap["$"+op] = val
        }

        out[field] = opMap

      default:
        return nil, fmt.Errorf("invalid filter value for '%s'", field)
    }
  }

  return out, nil
}

type QueryDSL struct {
  Filter  filter    `yaml:"Filter,omitempty"`
  Select  []string  `yaml:"Select,omitempty"`
  Sort    sort      `yaml:"Sort,omitempty"`
  Limit   int64     `yaml:"Limit,omitempty"`
}
