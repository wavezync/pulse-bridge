{{- define "pulse-bridge.name" -}}
{{- .Chart.Name -}}
{{- end -}}

{{- define "pulse-bridge.labels" -}}
app.kubernetes.io/name: {{ include "pulse-bridge.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{- define "pulse-bridge.selectorLabels" -}}
app.kubernetes.io/name: {{ include "pulse-bridge.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}