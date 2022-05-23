# go-netbox-proxy

## Authors

- [@twink0r - Alexander Karl](https://github.com/twink0r)
- [@Riege Software International GmbH](https://github.com/riege)

## State

Due a breaking changes in Netbox 3.2 we got some trouble with Terraform Providers. As this is opensource software the updates to fix this breaking change are dangling.

> **_breaking change:_**  The `created` field of all change-logged models now conveys a full datetime object, rather than only a date. (Previous date-only values will receive a timestamp of 00:00.) While this change is largely unconcerning, strictly-typed API consumers may need to be updated.
## Deployment

This small proxy is deployed between the `Ingress` (Varnish) and the `Service` (Netbox).
All requests to the `/api/` Endpoints will proxied by go-netbox-proxy. The content will scanned and substituted if the RegEX `"created":"(?:.+?)",` matches.

> **_NOTE:_**  This is a monkey patch and will not longer stay in place as needed!


To ensure that only API traffic will be forwarded through the proxy, this change is in place:
```
backend api {
  .host = "netbox-proxy";
}

sub vcl_recv {
  # Route API Traffic to the API Backend
  if (req.url ~ "^/api/") {
      set req.backend_hint = api;
  } else {
      set req.backend_hint = default;
  }
}
```
## Tech Stack

**Service:** HTTP, Proxy

**Cloud:** Kubernetes, OCI

## Feedback

If you have any feedback, please reach out to us at Riege Software Teams @ Core Channel
