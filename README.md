# Templates for ModelHelper

Templates for the [ModelHelper CLI](https://github.com/Modelhelper/cli) (`mh`). A template is a YAML file with a body written in the Go template language. `mh` fills the body with data from a database table (or a project name, or a file) and prints or writes the result.

| Folder | Contents |
|--------|----------|
| `code/` | Code templates, grouped by language (`go`, `cs`, `ts`, `sql`, `docker`, ...) |
| `db/mssql`, `db/postgres` | SQL building blocks, one set for each database. Code templates include them |
| `project/`, `projects/` | Project templates, used by `mh project create` |
| `solution/` | Solution templates |

## Contents

- [Quick start](#quick-start)
- [How templates are found and named](#how-templates-are-found-and-named)
- [The template file](#the-template-file)
- [Golang template language](#golang-template-language)
- [Functions added by ModelHelper](#functions-added-by-modelhelper)
- [Data input, or the model](#data-input---or-the-model)
- [Blocks: reuse part of a template](#blocks-reuse-part-of-a-template)
- [SQL for different databases](#sql-for-different-databases)
- [Snippets](#snippets)
- [Project templates](#project-templates)
- [Try and debug a template](#try-and-debug-a-template)

## Quick start

Install the CLI, then clone this repository and point `mh` at it. The first run of `mh` asks for the folders. They are saved in `~/.modelhelper/config.yaml`:

```yaml
templates:
  code:
    - /path/to/templates/code
  database:
    - /path/to/templates/db
  project:
    - /path/to/templates/projects
```

List the templates, then generate code. The demo source needs no database:

```bash
mh code template list --lang go
mh code generate -t go-struct --demo -s order --screen
```

With a real database, use `-s Schema.Table` instead of `--demo -s order`.

## How templates are found and named

`mh` reads every `.yaml` and `.yml` file below the folders in `templates.code`. The name of a template comes from its **path**, relative to that folder. Slashes become dashes and the extension is dropped:

| File | Name |
|------|------|
| `code/go/struct.yaml` | `go-struct` |
| `code/go/blocks/db/paged.yaml` | `go-blocks-db-paged` |
| `db/mssql/sql-select-all.yml` | `sql-select-all` |

A `name:` property inside the file is ignored, so the file path is the only thing to get right.

If two folders contain the same name, the later folder in `templates.code` wins. A project can also keep its own templates in `.modelhelper/templates/`. They are named the same way and replace global templates with the same name.

## The template file

```yaml
# mandatory, the template format version
version: 3

# mandatory, connects with the language definition (data types and so on)
language: go

# file    - produces a file, needs fileName if you want it exported
# block   - a part that other templates include, not generated on its own
# snippet - code that is inserted into an existing file
type: file

# optional, short text shown in 'mh code template list'
short: A struct for a table

# optional, longer text
description: |
  Creates a Go struct with one field for each column.

# which data the template receives, see "Data input" below
# basic (default) | name | entity | entities | custom | commits
model: entity

# used with 'mh code generate --export-bykey': the key picks the folder
# from the 'locations' in the project's .modelhelper/config.yaml
key: model

# the file name, a template itself. Used when the code is exported to disk
fileName: "{{ .Name | snake | singular }}.go"

# optional, groups this template is part of. 'mh code generate -f <group>'
# generates every template in the group
features:
  - go-rest-full

# optional, free labels
tags:
  - go
  - struct

# snippets only: the marker the code is inserted at, see "Snippets"
# identifier: models

# mandatory, the body, written in the Go template language
body: |
  package models

  type {{ .Name | singular | pascal }} struct {
  {{- range .Columns }}
    {{ .Name | pascal }} {{ .DataType | datatype }}
  {{- end }}
  }
```

## Golang template language

The body of a template is a standard Go [`text/template`](https://pkg.go.dev/text/template). The same language is used by Helm, Hugo and many other tools, so most guides on "Go templates" apply here. This section covers the parts you need to write ModelHelper templates. For the full reference, read the [`text/template` documentation](https://pkg.go.dev/text/template).

### Actions

Everything between `{{` and `}}` is an action. Everything outside is copied to the output as it is.

```
Hello {{ .Name }}      prints the Name field of the data
{{/* a comment */}}    prints nothing
```

### The dot

The dot (`.`) is the current data. At the top of a template, `.` is the model that `mh` hands to the template (see [Data input](#data-input---or-the-model)). `.Name` reads the field `Name` of the dot. Inside `range` and `with` the dot changes to the current item. Use `$` to get back to the model at the top:

```
{{ range .Columns }}
  {{ .Name }} belongs to {{ $.Name }}
{{ end }}
```

### Pipelines and functions

A pipeline passes a value through functions, left to right. Functions are called with their arguments after the name:

```
{{ .Name | singular | pascal }}      same as: pascal (singular .Name)
{{ printf "%s_id" .Name }}           a built-in function with arguments
{{ "ThisIsText" | snake }}           this_is_text
```

### Conditions

```
{{ if .IsPrimaryKey }}primary key{{ else if .IsNullable }}nullable{{ else }}required{{ end }}
{{ if and .HasChildren (not .HasParents) }}...{{ end }}
{{ if eq .Type "view" }}...{{ end }}
```

`false`, `0`, an empty string and an empty list count as false. The comparison functions are `eq`, `ne`, `lt`, `le`, `gt` and `ge`, combined with `and`, `or` and `not`.

### Loops

```
{{ range .Columns }}{{ .Name }},{{ end }}

{{ range $i, $col := .Columns }}{{ if $i }}, {{ end }}{{ $col.Name }}{{ end }}

{{ range .Children }}...{{ else }}no children{{ end }}
```

Columns have `IsFirst` and `IsLast`, which is often easier than counting:

```
{{ range .Columns }}{{ .Name }}{{ if not .IsLast }}, {{ end }}{{ end }}
```

### with

`with` runs its block only when a value is not empty, and moves the dot into it:

```
{{ with .Children }}{{ len . }} child tables{{ end }}
```

### Variables

```
{{ $model := .Name | camel | singular }}
{{ range .Columns }}{{ $model }}.{{ .Name | pascal }}{{ end }}
```

### White space

An action prints the white space around it. A dash removes it: `{{-` removes the white space before the action, `-}}` removes it after. Most templates use `{{-` on `range`, `if` and `end` lines so they do not leave blank lines:

```
type {{ .Name | pascal }} struct {
{{- range .Columns }}
  {{ .Name | pascal }} string
{{- end }}
}
```

### Include other templates

`{{ template "name" . }}` runs another template and prints the result. The last argument is the data it receives, usually `.`. See [Blocks](#blocks-reuse-part-of-a-template).

### Write a literal `{{`

The characters `{{` always start an action. To print them, print them as a string:

```
{{ "{{" }} .Name }}
```

The same goes for a `{` right before an action: `${{ .Name }}` is read as `$` and an action. Write `${{ "{" }}{{ .Name }}}` to get `${ORDER}`. Project templates are more forgiving, see [Project templates](#project-templates).

### Where to read more

- [`text/template` package](https://pkg.go.dev/text/template): the language reference
  - [Actions](https://pkg.go.dev/text/template#hdr-Actions), [Pipelines](https://pkg.go.dev/text/template#hdr-Pipelines), [Variables](https://pkg.go.dev/text/template#hdr-Variables) and [Functions](https://pkg.go.dev/text/template#hdr-Functions) (the built-in `eq`, `len`, `index`, `printf` and others)
  - [Text and spaces](https://pkg.go.dev/text/template#hdr-Text_and_spaces): the `{{-` and `-}}` trimming rules
- [Go by Example: Text Templates](https://gobyexample.com/text-templates): short, runnable examples
- [Go template syntax in Hugo](https://gohugo.io/templates/introduction/): another good explanation, but Hugo adds its own functions

## Functions added by ModelHelper

These are available in every code template, besides the built-in functions.

| Function | Example | Result |
|----------|---------|--------|
| `singular`, `plural` | `{{ "Child" \| plural }}` | `Children` |
| `lower`, `upper` | `{{ "Name" \| upper }}` | `NAME` |
| `camel` | `{{ "ThisIsASentence" \| camel }}` | `thisIsASentence` |
| `pascal` | `{{ "this_is_a_sentence" \| pascal }}` | `ThisIsASentence` |
| `snake` | `{{ "ThisIsASentence" \| snake }}` | `this_is_a_sentence` |
| `kebab` | `{{ "ThisIsASentence" \| kebab }}` | `this-is-a-sentence` |
| `macro` | `{{ "ThisIsASentence" \| macro }}` | `THIS_IS_A_SENTENCE` |
| `train` | `{{ "ThisIsASentence" \| train }}` | `This_Is_A_Sentence` |
| `dot` | `{{ "ThisIsASentence" \| dot }}` | `This.Is.A.Sentence` |
| `title` | `{{ "this is a sentence" \| title }}` | `This Is A Sentence` |
| `words` | `{{ "ThisIsASentence" \| words }}` | `This Is A Sentence` |
| `sentence` | `{{ "ThisIsASentence" \| sentence }}` | `This is a sentence` |
| `append` | `{{ "Order" \| append "Repository" }}` | `OrderRepository` |
| `increment` | `{{ increment $index }}` | `$index + 1`, handy for `$1, $2` parameters |
| `datatype` | `{{ .DataType \| datatype }}` | the language type for a database type |
| `nullable` | `{{ .DataType \| nullable true }}` | the language type in its nullable form (`sql.NullInt32` for `int` in Go). With `false` it is the plain type |
| `datatypeN` | same as `nullable` | |

`datatype` and `nullable` use the type mapping of the template's `language`, from the language definitions in your config (`languages.definitions`). A type that is not in the mapping is printed as it is. Run `mh language list` to see the languages and where each comes from.

`mh` has definitions for `cpp`, `cs`, `go`, `sql` and `ts` built in. Set `languages.definitions` in `~/.modelhelper/config.yaml` to a folder of `.yaml` files to add a language or to replace a built-in one, no rebuild needed. A file there with `language: go` replaces the built-in Go definition. Copy a built-in one as a start, and change the `datatypes`:

```yaml
language: rust
version: "3.0"
short: Language definition for Rust
datatypes:
  int:
    notNull: i32
    nullable: Option<i32>
```

The `tutorial-strings` template (`code/tutorial/strings.yaml`) prints all of these:

```bash
mh code generate -t tutorial-strings --demo --screen
```

## Data input - or the model

The data a template receives is the **model**. Which model you get is set with `model:` in the template.

| `model` | The template receives | Generated |
|---------|-----------------------|-----------|
| `basic` (default) | Project and namespace info | Once |
| `name` | The basic model and a `Name` from `-n` | Once |
| `entity` | One table, with columns and relations | Once for each table in `-s` |
| `entities` | A list of tables, in `.Entities` | Once, for all tables |
| `custom` | The content of a JSON file, in `.Custom` | Once. Use `--model file.json` |
| `commits` | Git commit history | Once. Used by the changelog templates |

```bash
# model: entity, one table: one set of code
mh code generate -t go-struct -s dbo.Order --screen

# model: entity, two tables: two sets of code
mh code generate -t go-struct -s dbo.Order -s dbo.Customer --screen

# model: entities, two tables: one set of code that loops over both
mh code generate -t go-structs-null -s dbo.Order -s dbo.Customer --screen

# model: name
mh code generate -t go-cli-cqrs-handler -n GetOrder --screen

# model: custom
mh code generate -t my-template --model data.json --screen
```

### The basic model

Every other model includes these fields. They are filled from the project's `.modelhelper/config.yaml`, if there is one (`setup`, `inject`, `options` and so on).

```go
type BasicModel struct {
	RootNamespace             string
	Namespace                 string            // setup.<key>.namespace
	Postfix                   string            // setup.<key>.postfix
	Prefix                    string            // setup.<key>.prefix
	ModuleLevelVariablePrefix string
	Inject                    []InjectSection
	Imports                   []string
	Project                   ProjectSection
	Developer                 DeveloperSection
	Options                   map[string]string // options in the project config
	PageHeader                string
	Feature                   FeatureModel
	Name                      string
}

type InjectSection struct {
	Name         string
	PropertyName string
}

type DeveloperSection struct {
	Name  string
	Email string
}

type ProjectSection struct {
	Name    string
	Owner   string
	Version string
	Exists  bool // true when the CLI runs in a project with a .modelhelper folder
}

// FeatureModel is filled from 'features' in the project config.
// UseLogger, UseTracing, UseMetrics, UseHealth, UseApi, UseDb, UseAuth, UseSwagger
// tell if a feature is switched on, Logger, Tracing, ... hold its options
type FeatureModel struct {
	UseLogger bool
	Logger    FeatureOptions
	// ...
}
```

### The entity model

With `model: entity` the template is run once for each table. It gets everything from the basic model, and these:

```go
type EntityModel struct {
	// ... all fields of the basic model, except Name

	ConnectionType    string // mssql, postgres, ...
	Name              string // the table or view name
	Schema            string
	Type              string // table | view
	Alias             string // built from the words in the name
	Synonym           string
	HasSynonym        bool
	ModelName         string // Synonym if there is one, otherwise Name
	Description       string
	HasDescription    bool
	HasPrefix         bool
	NameWithoutPrefix string

	Columns           []EntityColumnModel   // all visible columns
	NonPrimaryColumns []EntityColumnModel   // Columns without the primary key
	PrimaryKeys       []EntityColumnModel
	ForeignKeys       []EntityColumnModel
	UsedAsColumns     []EntityColumnModel   // created, modified, deleted markers and so on
	Parents           []EntityRelationModel // tables this table points to
	HasParents        bool
	Children          []EntityRelationModel // tables that point to this table
	HasChildren       bool

	UsesIdentityColumn bool
}

type EntityColumnModel struct {
	Name              string
	NameWithoutPrefix string
	DataType          string // the database type, use the datatype function
	Description       string
	Collation         string
	Length            int
	Precision         int
	Scale             int
	UseLength         bool
	UsePrecision      bool
	ReferencesTable   string
	ReferencesColumn  string

	IsPrimaryKey     bool
	IsForeignKey     bool
	IsIdentity       bool
	IsNullable       bool
	IsIgnored        bool
	IsDeletedMarker  bool
	IsCreatedDate    bool
	IsCreatedByUser  bool
	IsModifiedDate   bool
	IsModifiedByUser bool
	HasPrefix        bool
	HasDescription   bool

	IsFirst bool // first in the list it is part of
	IsLast  bool // last in the list it is part of
}

// Parents and Children
type EntityRelationModel struct {
	IsSelfJoin    bool
	RelatedColumn EntityColumnProps // the column on the other table
	OwnerColumn   EntityColumnProps // the column on this table
	OwnerName     string            // this table
	OwnerShcema   string

	// the other table, same fields as the entity: Name, Schema, Type, Alias, Synonym,
	// ModelName, Description, NameWithoutPrefix, Columns, PrimaryKeys, ...
	Name   string
	Schema string
}

type EntityColumnProps struct {
	Name       string
	DataType   string
	IsNullable bool
}
```

Run a template with `--demo` to see real values. The demo has the tables `order`, `orderItem`, `customer` and `product`.

### The entities model

`model: entities` gives the template the basic model and all tables in one list, so one file can contain code for several tables:

```go
type EntityListModel struct {
	// ... all fields of the basic model
	Entities []EntityModel
}
```

```
{{ range .Entities }}
  {{ .Name | pascal }}: {{ len .Columns }} columns
{{ end }}
```

## Blocks: reuse part of a template

A template with `type: block` is not generated on its own. Other templates include it with `template`, using its [name](#how-templates-are-found-and-named). The block gets the data you pass as the last argument, so a block for an entity can be used from any template that has an entity.

`code/go/blocks/column-list.yaml`:

```yaml
version: 3
language: go
type: block
short: Comma separated column names
model: entity
body: |
  {{- range .Columns }}{{ .Name }}{{ if not .IsLast }}, {{ end }}{{ end -}}
```

`code/go/hello-struct.yaml`:

```yaml
version: 3
language: go
type: file
key: model
model: entity
fileName: "{{ .Name | snake | singular }}.go"
body: |
  package models

  // {{ .Name | singular | pascal }} maps {{ .Schema }}.{{ .Name }}
  // columns: {{ template "go-blocks-column-list" . }}
  type {{ .Name | singular | pascal }} struct {
  {{- range .Columns }}
    {{ .Name | pascal }} {{ .DataType | datatype }} `json:"{{ .Name | camel }}"`{{ if .IsPrimaryKey }} // primary key{{ end }}
  {{- end }}
  }
```

```bash
mh code generate -t go-hello-struct --demo -s order --screen
```

```go
package models

// Order maps dbo.order
// columns: OrderId, CustomerId, Id, TotalAmount
type Order struct {
  OrderId string `json:"orderId"`
  CustomerId int `json:"customerId"`
  Id int `json:"id"` // primary key
  TotalAmount decimal `json:"totalAmount"`
}
```

(The built-in Go definition maps `decimal` to `decimal`, which is not a Go type. Override it with your own language definition, see [Functions added by ModelHelper](#functions-added-by-modelhelper).) Blocks can include other blocks. A template that names a block that does not exist fails with `template "x" not defined`.

## SQL for different databases

SQL that differs between databases lives in `db/mssql` and `db/postgres`, in files with the same names in both. Code templates include the SQL they need:

```
query := `{{ template "sql-select-all" . }}`
```

`mh` loads only the folder for the database type of the connection you use (`mssql` or `postgres`). The same Go template therefore produces `[dbo].[order]` and `@Id` for SQL Server, and `$1` for PostgreSQL. The demo source uses SQL Server. Use `mh code template list --db postgres` to list the templates for PostgreSQL.

If you add a SQL block, add it to **both** folders with the same file name, or the Go templates that include it fail for one database. The names used by the Go templates are:

| Name | Does |
|------|------|
| `sql-select-all` | All rows |
| `sql-select-single` | One row, by primary key |
| `sql-select-paged` | A page of rows, with the total row count |
| `sql-select-child-items` | The rows of a child table for a parent key |
| `sql-insert-item` | Insert and return the new id |
| `sql-update-item` | Update by primary key and return the key |
| `sql-delete-item` | Delete by primary key and return the id |

## Snippets

A template with `type: snippet` does not create a file. Its code is inserted into a file that already exists, at a marker. Put the marker `%%identifier%%` in the file, and set `identifier` in the template to the same word. The destination file comes from the `key` and `fileName`, as for other templates. Use `--export-bykey` to insert it.

## Project templates

Project templates in `project/` and `projects/` create a whole project with `mh project create`:

```bash
mh project template list
mh project create -t go-hex-api -n flight-api --dry-run
mh project create -t go-hex-api -n flight-api
```

A project template is a folder with a `project.yaml` and the files to copy:

```yaml
templateVersion: 3.0
name: go-hex-api
language: go
description: Creates an API in Go
version: 0.0.1
# the name of the folder to create the project in, empty means the current folder.
# It is a template too, for example "{{ .Name | kebab }}"
rootDirectory: ""
# the folders to copy, relative to the project.yaml file
sources:
  - ./src
```

The folder structure is copied as it is, including empty folders. Every file is run through the Go template engine, with this data:

| Field | Value |
|-------|-------|
| `.Name` | The name given with `-n` |
| `.Version` | The version from `project.yaml`, or `0.0.1` |

These functions are available: `lower`, `upper`, `words`, `sentence`, `snake`, `macro`, `train`, `kebab`, `dot`, `title`, `proper` (same as `title`), `pascal` and `camel`. `plural` and `singular` return the text unchanged here.

Some things to know:

- A file that **cannot be parsed** is copied as it is, and `mh` prints `Warning: '<file>' was copied as is, it could not be parsed`. If you see untouched `{{ ... }}` in a created project, look for that warning.
- `${{{ .Name | snake | upper }}_VERSION}` works. It prints `${MY_API_VERSION}`, which is handy in docker and shell files.
- To keep a `{{ ... }}` as it is in a created file, write it as a string: `{{ "{{ .Name }}" }}`. This is needed in the `.modelhelper/config.yaml` of a project template, where `locations` contain templates of their own.
- Put a `.modelhelper` folder with a `config.yaml` and a `README.md` inside `src` to give the new project its code generation setup. See the [CLI documentation](https://github.com/Modelhelper/cli#using-it-in-a-project).

## Try and debug a template

```bash
# the templates mh can see, with their keys and models
mh code template list --lang go --db mssql

# one template: the key, file name, where it is from, and the body
mh code template detail go-struct

# run a template on the built-in sample tables, no database needed
mh code generate -t go-struct --demo -s order --screen

# the same, with only the code in the output
mh code generate -t go-struct --demo -s order --screen --code-only
```

If a template fails, `mh` prints `Error when generating template '<name>': ...` with the line and column in the template. The most common causes:

- `function "x" not defined`: a typo, or a function from another tool. See [the list](#functions-added-by-modelhelper).
- `template "x" not defined`: the name in `{{ template "x" . }}` does not match a template. Check the [name rules](#how-templates-are-found-and-named), and for a SQL block that it exists in the database folder you use.
- `can't evaluate field X in type ...`: the model has no such field. Check the [model](#data-input---or-the-model).
- `wrong number of args`: for example `nullable` takes two arguments, `{{ .DataType | nullable true }}`.

For Go output, save the code and run `gofmt -e` on it to catch syntax errors in what you generate.

Agent skills that help write and check templates, including a script that finds missing `template` names, are in [Modelhelper/skills](https://github.com/Modelhelper/skills).
