# 0.1.7 (2023-02-15)
 - Fix regression in 0.1.6 whereby the frontend credentials were always created

# 0.1.6 (2023-02-14)
 - Fix issue that defaults boolean values on stores and taxes aren't set

# 0.1.5 (2023-02-02)
 - Fix newline error in template when rendering store variables
 - Remove duplicate `enabled` property in the `messages` block

# 0.1.4 (2023-02-02)
 - Fix template error when using managed stores

# 0.1.3 (2023-02-01)
 - Refactor the hcl renderer to use go templates. This should also fix escaping
   issues for some values.

# 0.1.2 (2023-02-01)
 - Fix goreleaser file to be compatible with plugin registry

# 0.1.1 (2023-01-31)
 - Fix invalid rendering of store variables in the generated terraform files.
 - Fix explicitly disabling the creation of frontend credentials.
 - Fix setting store_secrets and store_variablds in the commercetools block

# 0.1.0 (2023-01-18)
Initial creation of the changelog
