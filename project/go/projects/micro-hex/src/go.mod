module {{ .Name | kebab }}

go {{ .Version.Go }}

{{ if .HasPackages }}
require (
	{{ range .Packages }}{{ .Name }} {{ .Version }} {{ end }}
)
{{ end }}