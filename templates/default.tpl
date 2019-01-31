{{ range $key, $value := .Fields }}
   <b>{{ $key }}</b>: {{ $value }}
{{ end }}
