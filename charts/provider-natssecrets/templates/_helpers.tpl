{{/*
Expand the name of the chart.
*/}}
{{- define "provider-natssecrets.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "provider-natssecrets.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "provider-natssecrets.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "provider-natssecrets.labels" -}}
helm.sh/chart: {{ include "provider-natssecrets.chart" . }}
{{ include "provider-natssecrets.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "provider-natssecrets.selectorLabels" -}}
app.kubernetes.io/name: {{ include "provider-natssecrets.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

# {{/*
# Create the name of the service account to use
# */}}
# {{- define "provider-natssecrets.serviceAccountName" -}}
# {{- if .Values.serviceAccount.create }}
# {{- default (include "provider-natssecrets.fullname" .) .Values.serviceAccount.name }}
# {{- else }}
# {{- default "default" .Values.serviceAccount.name }}
# {{- end }}
# {{- end }}

{{- define "vaultData" -}}
{{- $vault := .Values.vault.token.secretRef.data -}}
{{- required "vault.token.secretRef.data.address is required" $vault.address -}}
{{- required "vault.token.secretRef.data.tls is required" $vault.tls -}}
{{- required "vault.token.secretRef.data.insecure is required" $vault.insecure -}}
{{- required "vault.token.secretRef.data.token is required" $vault.token -}}
{
  "address": "{{ $vault.address }}",
  "tls": {{ $vault.tls }},
  "insecure": {{ $vault.insecure }},
  "token": "{{ $vault.token }}",
  "path": "nats-secrets"
}
{{- end -}}