# Contributing to the ModelHelper templates

Thanks for wanting to help. You can fix a template that produces wrong code, add a template for a language or a pattern, or improve the docs. Please read this page before you open a pull request.

- [Code of Conduct](./CODE_OF_CONDUCT.md)
- [Reporting issues](#reporting-issues)
- [Pull requests](#pull-requests)
- [Writing and changing templates](#writing-and-changing-templates)
- [Testing your change](#testing-your-change)
- [Committing changes](#committing-changes)
- [Repository structure](#repository-structure)
- [License](#license)

## Reporting issues

Use [GitHub issues](https://github.com/Modelhelper/templates/issues/new/choose). For a template that produces wrong code, include:

- the template name (`mh code template list`)
- the `mh code generate` command you ran
- the table definition, or the `--demo` table that shows the problem
- the output you got and the output you expected
- the database type (SQL Server or PostgreSQL)

Problems in the `mh` program itself belong in the [CLI repository](https://github.com/Modelhelper/cli/issues).

## Pull requests

- Branch from `main` and open the pull request against `main`.
- For a new template, open an issue first. Say which language and pattern it covers, so we can check that it does not overlap with an existing one.
- Keep a pull request to one change. Several small commits are fine, they are squashed when merged.
- Show the output. A short before/after of the generated code in the pull request description is the fastest way to a review.
- Update [README.md](../README.md) when you add something that changes how templates are written or found, such as a new model, function or folder.

## Writing and changing templates

The [README](../README.md) explains the file format, the Go template language, the functions and the models. These rules apply to this repository:

- **The path is the name.** `code/go/blocks/db/paged.yaml` is `go-blocks-db-paged`. The `name:` property is ignored. Renaming or moving a file renames the template, so search for `{{ template "old-name"` before you do it.
- **Follow the folder pattern.** Look at the templates next to yours. Put code templates in `code/<language>/`, reusable parts in a `blocks/` folder, and project templates in `projects/`.
- **SQL comes in pairs.** A SQL block in `db/mssql` needs a block with the same file name in `db/postgres`, in that database's syntax (`@name` and `[brackets]` for SQL Server, `$1` and `RETURNING` for PostgreSQL). A Go template that includes `{{ template "sql-select-all" . }}` fails for the database that is missing the block.
- **Set `short`.** It is shown in `mh code template list`.
- **Set `key` and `fileName` for templates that create files.** The key decides which project folder `--export-bykey` writes to, so use an existing key when one fits (`mh code template list` shows them).
- **Generated code must compile.** Imports, error handling and naming should be what you would accept in a code review of your own project.
- **Use the functions for naming.** `{{ .Name | singular | pascal }}`, not a hard-coded name.
- **Handle empty lists and the last item.** A table can have no children, and column lists need to leave out the trailing comma (`{{ if not .IsLast }}, {{ end }}`).

## Testing your change

You need the [`mh` CLI](https://github.com/Modelhelper/cli#install) and this repository in your `templates` config.

```bash
# the template is found and has the key and model you expect
mh code template list --lang go --db mssql

# render it on the built-in sample tables, no database needed
mh code generate -t go-struct --demo -s order --screen

# try more than one table: order has children and a foreign key, customer does not
mh code generate -t go-struct --demo -s customer --screen
```

Before you open the pull request, check these:

1. The template renders without an `Error when generating template` message.
2. It works for every sample table you can use (`order`, `orderItem`, `customer`, `product`), not only the first.
3. For SQL, render a Go template that includes your block with the SQL Server and the PostgreSQL connection, or at least read both outputs.
4. Generated Go code passes `gofmt -e`. Use `--code-only` and save the output to a file.
5. Every `{{ template "..." }}` name you use exists. The `create-modelhelper-template` skill in [Modelhelper/skills](https://github.com/Modelhelper/skills) has a script, `scripts/check-refs.py`, that lists the ones that do not. Run it from the root of this repository.

## Committing changes

Commit messages follow the [commit message convention](./COMMIT_CONVENTION.md), for example `fix(go): add missing slog import in repository` or `feat(sql): add paged query for postgres`. The convention is not enforced by a hook, so check your messages before you push.

## Repository structure

| Folder | Contents |
|--------|----------|
| `code/` | Code templates, one folder for each language, and shared blocks |
| `db/mssql`, `db/postgres` | SQL blocks, the same file names in both |
| `project/`, `projects/` | Project templates, used by `mh project create` |
| `solution/` | Solution templates |

## License

This repository is licensed under the [GNU General Public License v3.0](../LICENSE). By opening a pull request you agree that your contribution is licensed under the same license.
