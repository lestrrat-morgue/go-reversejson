go-reversejson
==============

This is a silly little tool to convert JSON into pseudo Go structures.
It's intended to be used for when you have a spec before hand with
giant JSONs, and you feel silly re-typing all those fields in as you
CamelCase the field names and add types to it, etc.

I don't know if this is really useful. But I just didn't feel like
cut n pasting stuff from JSON to Go.

As of this writing, you just do

```
cat foo.json | go run main.go -name StructName
```

and it outputs something that's close enough.

My use case is this: I was looking to implement some OpenID Connect stuff in go,
and found JSON structures like this (https://openid.net/specs/openid-connect-discovery-1_0.html#rfc.section.4.2):

```json
{
   "issuer":
     "https://server.example.com",
   "authorization_endpoint":
     "https://server.example.com/connect/authorize",
   "token_endpoint":
     "https://server.example.com/connect/token",
   "token_endpoint_auth_methods_supported":
     ["client_secret_basic", "private_key_jwt"],
   "token_endpoint_auth_signing_alg_values_supported":
     ["RS256", "ES256"],
   "userinfo_endpoint":
     "https://server.example.com/connect/userinfo",
   "check_session_iframe":
     "https://server.example.com/connect/check_session",
   "end_session_endpoint":
     "https://server.example.com/connect/end_session",
   "jwks_uri":
     "https://server.example.com/jwks.json",
   "registration_endpoint":
     "https://server.example.com/connect/register",
   "scopes_supported":
     ["openid", "profile", "email", "address",
      "phone", "offline_access"],
   "response_types_supported":
     ["code", "code id_token", "id_token", "token id_token"],
   "acr_values_supported":
     ["urn:mace:incommon:iap:silver",
      "urn:mace:incommon:iap:bronze"],
   "subject_types_supported":
     ["public", "pairwise"],
   "userinfo_signing_alg_values_supported":
     ["RS256", "ES256", "HS256"],
   "userinfo_encryption_alg_values_supported":
     ["RSA1_5", "A128KW"],
   "userinfo_encryption_enc_values_supported":
     ["A128CBC-HS256", "A128GCM"],
   "id_token_signing_alg_values_supported":
     ["RS256", "ES256", "HS256"],
   "id_token_encryption_alg_values_supported":
     ["RSA1_5", "A128KW"],
   "id_token_encryption_enc_values_supported":
     ["A128CBC-HS256", "A128GCM"],
   "request_object_signing_alg_values_supported":
     ["none", "RS256", "ES256"],
   "display_values_supported":
     ["page", "popup"],
   "claim_types_supported":
     ["normal", "distributed"],
   "claims_supported":
     ["sub", "iss", "auth_time", "acr",
      "name", "given_name", "family_name", "nickname",
      "profile", "picture", "website",
      "email", "email_verified", "locale", "zoneinfo",
      "http://example.info/claims/groups"],
   "claims_parameter_supported":
     true,
   "service_documentation":
     "http://server.example.com/connect/service_documentation.html",
   "ui_locales_supported":
     ["en-US", "en-GB", "en-CA", "fr-FR", "fr-CA"]
  }
```

So then I wrote this code, copied it in a local file, and ran

```
cat foo.json | go run main -name ProviderConfig
```

And out comes this (gofmt has been applied to output for better visibility):

```go
type ProviderConfig struct {
  AcrValuesSupported                         []string `json:"acr_values_supported"`
  AuthorizationEndpoint                      string   `json:"authorization_endpoint"`
  CheckSessionIframe                         string   `json:"check_session_iframe"`
  ClaimTypesSupported                        []string `json:"claim_types_supported"`
  ClaimsParameterSupported                   bool     `json:"claims_parameter_supported"`
  ClaimsSupported                            []string `json:"claims_supported"`
  DisplayValuesSupported                     []string `json:"display_values_supported"`
  EndSessionEndpoint                         string   `json:"end_session_endpoint"`
  IdTokenEncryptionAlgValuesSupported        []string `json:"id_token_encryption_alg_values_supported"`
  IdTokenEncryptionEncValuesSupported        []string `json:"id_token_encryption_enc_values_supported"`
  IdTokenSigningAlgValuesSupported           []string `json:"id_token_signing_alg_values_supported"`
  Issuer                                     string   `json:"issuer"`
  JwksUri                                    string   `json:"jwks_uri"`
  RegistrationEndpoint                       string   `json:"registration_endpoint"`
  RequestObjectSigningAlgValuesSupported     []string `json:"request_object_signing_alg_values_supported"`
  ResponseTypesSupported                     []string `json:"response_types_supported"`
  ScopesSupported                            []string `json:"scopes_supported"`
  ServiceDocumentation                       string   `json:"service_documentation"`
  SubjectTypesSupported                      []string `json:"subject_types_supported"`
  TokenEndpoint                              string   `json:"token_endpoint"`
  TokenEndpointAuthMethodsSupported          []string `json:"token_endpoint_auth_methods_supported"`
  TokenEndpointAuthSigningAlgValuesSupported []string `json:"token_endpoint_auth_signing_alg_values_supported"`
  UiLocalesSupported                         []string `json:"ui_locales_supported"`
  UserinfoEncryptionAlgValuesSupported       []string `json:"userinfo_encryption_alg_values_supported"`
  UserinfoEncryptionEncValuesSupported       []string `json:"userinfo_encryption_enc_values_supported"`
  UserinfoEndpoint                           string   `json:"userinfo_endpoint"`
  UserinfoSigningAlgValuesSupported          []string `json:"userinfo_signing_alg_values_supported"`
}
```

Currently it understands `string`, `float`, `bool`, and lists, but lists
are assumed to be list of strings.