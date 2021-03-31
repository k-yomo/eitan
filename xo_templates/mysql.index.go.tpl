{{- $short := (shortname .Type.Name "err" "sqlstr" "db" "q" "res" "XOLog" .Fields) -}}
{{- $table := (.Type.Table.TableName) -}}
// {{ .FuncName }} retrieves a row from '{{ $table }}' as a {{ .Type.Name }}.
// Generated from index '{{ .Index.IndexName }}'.
func {{ .FuncName }}(ctx context.Context, db Queryer{{ goparamlist .Fields true true }}) ({{ if not .Index.IsUnique }}[]{{ end }}*{{ .Type.Name }}, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`{{ colnames .Type.Fields }} ` +
		`FROM {{ $table }} ` +
		`WHERE {{ colnamesquery .Fields " AND " }}`

	// log and trace
	XOLog(ctx, sqlstr{{ goparamlist .Fields true false }})
	closeSpan := startSQLSpan(ctx, "{{ .FuncName }}", sqlstr{{ goparamlist .Fields true false }})
	defer closeSpan()

{{- if .Index.IsUnique }}
	{{ $short }} := {{ .Type.Name }}{
	{{- if .Type.PrimaryKey }}
		_exists: true,
	{{ end -}}
	}

	// run query
	err = db.QueryRowxContext(ctx, sqlstr{{ goparamlist .Fields true false }}).Scan({{ fieldnames .Type.Fields (print "&" $short) }})
	if err != nil {
		return nil, err
	}

	return &{{ $short }}, nil
{{- else }}
	// run query
	rows, err := db.QueryContext(ctx, sqlstr{{ goparamlist .Fields true false }})
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// load results
	var res []*{{ .Type.Name }}
	for rows.Next() {
		{{ $short }} := {{ .Type.Name }}{
		{{- if .Type.PrimaryKey }}
			_exists: true,
		{{ end -}}
		}

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

