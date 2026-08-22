{{- define "mini-erp.name" -}}
mini-erp
{{- end -}}

{{- define "mini-erp.fullname" -}}
{{ .Release.Name }}
{{- end -}}

{{- define "mini-erp.image" -}}
{{- printf "%s/%s/%s:%s" .root.Values.global.imageRegistry .root.Values.global.imageOwner .repository .root.Values.global.imageTag -}}
{{- end -}}

{{- define "mini-erp.labels" -}}
app.kubernetes.io/name: {{ include "mini-erp.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}
