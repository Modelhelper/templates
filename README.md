# templates for ModelHelper

## The yaml file type

All templates is defined in a yaml file. This is for improving the writeability and readability for the template author as well as making it easy for the application to read, use and structure the templates.

### Template structure

A "valid" template can have these properties

```yaml

# mandatory, this is the template version
version: 3.0 

# mandatory field. connects with correct language definition
language: cs 

# mandatory, default: file 
type: file | snippet | block

# optional
description: |
    Long description if needed. Will show when 'mh template <template>'

# optional
short: This should be short and to the point

# optional, default is null
tags: 
 - tag1
 - tag2
 - tag3 

# what groups this template is included in
# is used to fetch templates within a group when using -tg in generate command
# optional, default is null
groups: 
 - group1
 - group2
 - group3

# optional, default is none
model: none # values none | basic | entity | entities
# planned but not implemented models
## template | none | graphql | json | database | project | options 

# optional, settings here will only be used if --export or --export-file has been used as option argument in the cli
export:

  # if snippet is used, this will be the filename to inject the snippets to
  # optional 
  fileName: filename to be used in the export

  # optional. A key cannot contain . (dot)
  key: the key to be used to find correct path

# mandatory, the body of the template
body: |
    This is the template body written with golang text/template specification
    the {{ . }} notation corresponds to the import type.

    the . inside a {{ . }} provides the context of the given import

    e.g if the template import entities you can loop through all entities like so:
    the context of . in the range is a list of entities
        {{ range . }}
            inside the loop, the context of . is each entity
            {{ .Schema }} // the schema of the current entity
            {{ .Name }} // entity name
            {{ .RowCount }}
            {{ .Type }} // table | view | stored procedure | function...
            {{ .Created }} // date for creation
        {{ end }} // end of loop
```

## Golang template language

## Data input - or the model

Data pushed to a template is called a model. For a modelhelper template there are four types of models if you include the empty model.

### none (the empty model)

### The 'basic' model

#### the basic model usage in the modelhelper CLI

```bash
# generate code based on an 'entity' model template and only one table
# this will produce one set of code for that template
mh generate --template tutorial-model-entity -entity NameOfTable

# generate code based on an 'entity' model template and only one table
# Because we use two entities this command will produce two sets of code for that template
mh generate --template tutorial-model-entity -entity NameOfTable -entity AnotherTable
```

#### the basic model structure

```go
type BasicModel struct {
 RootNamespace             string
 Namespace                 string
 Postfix                   string
 Prefix                    string
 ModuleLevelVariablePrefix string
 Inject                    []InjectSection
 Imports                   []string
 Project                   ProjectSection
 Developer                 DeveloperSection
 Options                   map[string]string
 PageHeader                string
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
 Name  string
 Owner string
}
```

### The entity model

When the ```model``` property of a template is set to ```entity``` each template is served with exactely one entity at the time. This means that you can create one file for each ```entity```.
You can still enter as many ```entities``` from the modelhelper CLI as you wish, but the template works with one entity and then move to the next.

#### entity model usage in the modelhelper CLI

```bash
# generate code based on an 'entity' model template and only one table
# this will produce one set of code for that template
mh generate --template tutorial-model-entity -entity NameOfTable

# generate code based on an 'entity' model template and only one table
# Because we use two entities this command will produce two sets of code for that template
mh generate --template tutorial-model-entity -entity NameOfTable -entity AnotherTable
```

#### entity model structure

Use this structure of the entity model as a guide when you write templates for your self or a modelhelper core template

The entity model includes all the properties and values from the BasicModel

```go
// EntityModel
type EntityModel struct {
 // properties from the basic model
 // this model is explained in details above this secion
 RootNamespace             string
 Namespace                 string
 Postfix                   string
 Prefix                    string
 ModuleLevelVariablePrefix string
 Inject                    []InjectSection
 Imports                   []string
 Project                   ProjectSection
 Developer                 DeveloperSection
 Options                   map[string]string
 PageHeader                string

 // specific properties for the entity model
 Name                      string   //this is the name of the source (tablename, view name)
 Schema                    string   //the database schema the entity is on
 Type                      string   //the type of object (Table, View)
 Alias                     string   //Alias built on the words of the table name)
 Synonym                   string   //If indicated in the source 
 HasSynonym                bool     //If a synonym exists
 ModelName                 string   //A merged value of Synonym or Name
 Description               string   //The description of the entity
 HasDescription            bool     //If the source has a description
 HasPrefix                 bool     //If there is a prefix on this entity name
 NameWithoutPrefix         string   //The entity name without the prefix

 Columns                   []EntityColumnModel      // A list of all visible columns
 Parents                   []EntityRelationModel    // A list of all 'Parent' relations
 Children                  []EntityRelationModel    // A list of all 'Child' relations
 PrimaryKeys               []EntityColumnModel      // A list of all primary keys
 ForeignKeys               []EntityColumnModel      // A list of all foreign keys
 UsedAsColumns             []EntityColumnModel      // A list of all special columns like modified, changed, created date and so on
 UsesIdentityColumn        bool     // Indicates if the table uses an identity column
}

type EntityColumnModel struct {
 Description       string
 IsForeignKey      bool
 IsPrimaryKey      bool
 IsIdentity        bool
 IsNullable        bool
 IsIgnored         bool
 IsDeletedMarker   bool
 IsCreatedDate     bool
 IsCreatedByUser   bool
 IsModifiedDate    bool
 IsModifiedByUser  bool
 HasPrefix         bool
 HasDescription    bool
 Name              string
 NameWithoutPrefix string
 Collation         string
 ReferencesColumn  string
 ReferencesTable   string
 DataType          string
 Length            int
 Precision         int
 Scale             int
 UseLength         bool
 UsePrecision      bool
}

type EntityRelationModel struct {
 IsSelfJoin        bool
 ReleatedColumn    EntityColumnProps 
 OwnerColumn       EntityColumnProps
 Name              string
 Schema            string
 Type              string
 Alias             string
 Synonym           string
 HasSynonym        bool
 ModelName         string
 Description       string
 HasDescription    bool
 HasPrefix         bool
 NameWithoutPrefix string
 UsesIdentityColumn bool
}

type EntityColumnProps struct {
 Name       string
 DataType   string
 IsNullable bool
}
```

### The entities model

With the entities model you can iterate over the entities you use as the input source. Meaning one code template file can contains data about multiple input sources (tables).

#### entities model structure

The entities model includes all the properties and values from the BasicModel

```go
type EntityListModel struct {
 RootNamespace             string
 Namespace                 string
 Postfix                   string
 Prefix                    string
 ModuleLevelVariablePrefix string
 Inject                    []InjectSection
 Imports                   []string
 Project                   ProjectSection
 Developer                 DeveloperSection
 Options                   map[string]string
 PageHeader                string
 Entities                  []EntityModel
}
```

#### entities model usage in the modelhelper CLI
