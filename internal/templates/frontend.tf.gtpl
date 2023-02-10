{{ if .Config.Stores }}
{{ range $store := .Config.Stores }}
resource "commercetools_api_client" "frontend_credentials_{{ $store.Key }}" {
  name = "frontend_credentials_terraform_{{ $store.Key }}"
  scope = {{ RenderScopes $.Config.Frontend.PermissionScopes $.Config.ProjectKey $store.Key }}

  {{ if derefBool $store.Managed true }}
  depends_on = [commercetools_store.{{ $store.Key }}]
  {{ end }}
}

output "frontend_client_scope_{{ $store.Key }}" {
  value = commercetools_api_client.frontend_credentials_{{ $store.Key }}.scope
}

output "frontend_client_id_{{ $store.Key }}" {
  value = commercetools_api_client.frontend_credentials_{{ $store.Key }}.id
}

output "frontend_client_secret_{{ $store.Key }}" {
  value = commercetools_api_client.frontend_credentials_{{ $store.Key }}.secret
}
{{ end }}
{{ else }}
resource "commercetools_api_client" "frontend_credentials" {
  name = "frontend_credentials_terraform"
  scope = {{ RenderScopes .Config.Frontend.PermissionScopes .Config.ProjectKey "" }}
}

output "frontend_client_scope" {
    value = commercetools_api_client.frontend_credentials.scope
}

output "frontend_client_id" {
    value = commercetools_api_client.frontend_credentials.id
}

output "frontend_client_secret" {
    value = commercetools_api_client.frontend_credentials.secret
}
{{ end }}
