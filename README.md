# Charlatan

[![Build Status](https://travis-ci.org/BatchLabs/charlatan.svg?branch=master)](https://travis-ci.org/BatchLabs/charlatan)

**Charlatan** is a query engine for lists or streams of records in different
formats. It natively supports CSV and JSON formats but can easily be extended
to others.

It supports an SQL-like query language that is defined below. Queries are
applied to records to extract values depending on zero or more criteria.

## Query Syntax

```
SELECT <fields> FROM <source> [ WHERE <value> ] [ STARTING AT <index> ]
```

- `<fields>` is a list of comma-separated field names. Each field name must
  exist in the source. When reading CSV files, the field names are the column
  names, while when reading JSON they represent keys.
- `<source>` is the filename from which the data is read. The API is agnostique
  on this and one can implement support for any source type.
- `<value>` is a SQL-like value, which can be either a constant (e.g.
  `WHERE 1`), a field (e.g. `WHERE archived`) or any operation using comparison
  operators (`=`, `!=`, `<`, `<=`, `>`, `>=`, `AND`, `OR`) and optionally
  parentheses (e.g. `WHERE (foo > 2) AND (bar = "yo")`). The parser allows to
  use `&&` instead of `AND` and `||` instead of `OR`. It also support inclusive
  range tests, like `WHERE age IN [20, 30]`.
- `<index>` can be used to skip the first N records.

Constant values include strings, integers, floats, booleans and the `null`
value.

### Examples

```sql
SELECT CountryName FROM sample/csv/population.csv WHERE Year = 2010 AND Value > 50000000 AND Value < 70000000
SELECT name, age FROM sample/json/people.jsons WHERE stats.walking > 30 AND stats.biking < 300
SELECT name, age FROM sample/json/people.jsons WHERE stats.walking IN [20, 100]
```

## API

The library is responsible for parsing the query and executing against records.
Everything else is up to you, including how fields are retrieved from records.

Note: code examples below don’t include error handling for clarity purposes.

```go
// parse the query
query, _ := charlatan.QueryFromString("SELECT foo FROM myfile.json WHERE foo > 2")

// open the source file
reader, _ := os.Open(query.From())

defer reader.Close()

// skip lines if the query contains "STARTING AT <n>"
skip := query.StartingAt()

decoder := json.NewDecoder(reader)

for {
    skip--
    if skip >= 0 {
        continue
    }

    // get a new JSON record
    r, err := record.NewJSONRecordFromDecoder(decoder)

    if err == io.EOF {
        break
    }

    // evaluate the query against the record to test if it matches
    if match, _ := query.Evaluate(r); !match {
        continue
    }

    // extract the values and print them
    values, _ := query.FieldsValues(r)
    fmt.Printf("%v\n", values)
}
```

Two record types are included: `JSONRecord` and `CSVRecord`. Implementing a
record only requires one method: `Find(*Field) (*Const, error)`, which takes a
field and return its value.

As an example, let’s implement a `LineRecord` that’ll be used to get specific
characters on each line of a file, `c0` being the first character:

```go
type LineRecord struct { Line string }

func (r *LineRecord) Find(f *charlatan.Field) (*charlatan.Const, error) {

    // this is the field value we must return
    name := f.Name()

    // we reject fields that doesn't start with 'c'
    if len(name) < 2 || name[0] != 'c' {
        return nil, fmt.Errorf("Unknown field '%s'", name)
    }

    // we extract the character index from the field name.
    index, err := strconv.ParseInt(name[1:], 10, 64)
    if err != nil {
        return nil, err
    }

    // let's not be too strict and accept out-of-range indexes
    if index < 0 || index >= int64(len(r.Line)) {
        return charlatan.StringConst(""), nil
    }

    return charlatan.StringConst(fmt.Sprintf("%c", r.Line[index])), nil
}
```

One can now loop over a file’s content, construct `LineRecord`s from its lines
and evaluate queries against them:

```go
query, _ := charlatan.QueryFromString("SELECT c1 FROM myfile WHERE c0 = 'a'")

f, _ := os.Open(query.From())
defer f.Close()

s := bufio.NewScanner(f)
for s.Scan() {
    r := &LineRecord{Line: s.Text()}

    if m, _ := query.Evaluate(r); !m {
        continue
    }

    values, _ := query.FieldsValues(r)
    fmt.Printf("%v\n", values)
}
```

### Examples

Two examples are included in the repository under `sample/csv/` and
`sample/json/`.

### Authors

- [Nicolas DOUILLET](https://github.com/minimarcel)
- [Vincent RISCHMANN](https://github.com/vrischmann)
- [Baptiste FONTAINE](https://github.com/bfontaine)

