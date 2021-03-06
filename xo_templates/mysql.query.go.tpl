{{- $short := (shortname .Type.Name "err" "sqlstr" "db" "q" "res" "XOLog" .QueryParams) -}}
{{- $queryComments := .QueryComments -}}
{{- if .Comment -}}
// {{ .Comment }}
{{- else -}}
// {{ .Name }} runs a custom query, returning results as {{ .Type.Name }}.
{{- end }}
func {{ .Name }} (ctx context.Context, db Queryer{{ range .QueryParams }}, {{ .Name }} {{ .Type }}{{ end }}) ({{ if not .OnlyOne }}[]{{ end }}*{{ .Type.Name }}, error) {
	var err error

	// sql query
	{{ if .Interpolate }}var{{ else }}const{{ end }} sqlstr = {{ range $i, $l := .Query }}{{ if $i }} +{{ end }}{{ if (index $queryComments $i) }} // {{ index $queryComments $i }}{{ end }}{{ if $i }}
	{{end -}}`{{ $l }}`{{ end }}

    // log and trace
	XOLog(ctx, sqlstr{{ range .QueryParams }}{{ if not .Interpolate }}, {{ .Name }}{{ end }}{{ end }})
	closeSpan := startSQLSpan(ctx, {{ .Name }}, sqlstr{{ range .QueryParams }}{{ if not .Interpolate }}, {{ .Name }}{{ end }}{{ end }})
	defer closeSpan()

	// run query
{{- if .OnlyOne }}
	var {{ $short }} {{ .Type.Name }}
	err = db.QueryRowContext(ctx, sqlstr{{ range .QueryParams }}, {{ .Name }}{{ end }}).Scan({{ fieldnames .Type.Fields (print "&" $short) }})
	if err != nil {
		return nil, err
	}

	return &{{ $short }}, nil
{{- else }}
	rows, err := db.QueryContext(ctx, sqlstr{{ range .QueryParams }}, {{ .Name }}{{ end }})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// load results
	res := []*{{ .Type.Name }}{}
	for rows.Next() {
		{{ $short }} := {{ .Type.Name }}{}

		// scan
		err = rows.Scan({{ fieldnames .Type.Fields (print "&" $short) }})
		if err != nil {
			return nil, err
		}

		res = append(res, &{{ $short }})
	}

	return res, nil
{{- end }}
}

