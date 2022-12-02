# Commercetools Plugin for Mach Composer 

This repository contains the Sentry plugin for Mach Composer. It requires Mach Composer 3.x


## Usage

```yaml
global:
  # ...

sites:
  - identifier: my-site

    commercetools:
      project_key: project-key
      client_id: client-id
      client_secret: client-secret
      scopes: manage_project:project-key manage_api_clients:project-key view_api_clients:project-key
      
```
