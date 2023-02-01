{{ range $store := .Config.ManagedStores }}
resource "commercetools_store" "{{ $store.Key }}" {
  key  = "{{ $store.Key }}"

  {{ renderProperty "name" $store.Name }}
  {{ renderProperty "languages" $store.Languages }}

  {{ if $store.DistributionChannels }}
  distribution_channels = [
    {{ range $c := $store.DistributionChannels -}}
      commercetools_channel.{{ $c }}.key,
    {{- end }}
  ]
  {{ end }}
  {{ if $store.SupplyChannels }}
  supply_channels = [
    {{ range $c := $store.SupplyChannels -}}
      commercetools_channel.{{ $c }}.key,
    {{- end }}
  ]
  {{ end }}
}
{{ end }}
