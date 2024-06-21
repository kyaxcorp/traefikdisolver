# Real IP from Cloudflare/AWS Cloudfront Proxy/Tunnel

If Traefik is behind a Cloudflare/AWS Cloudfront Proxy/Tunnel, it won't be able to get the real IP from the external client as well as other information.

Processed Headers:
- Cloudflare: CF-Connecting-IP
- Cloudfront: Cloudfront-Viewer-Address

## Configuration

### Configuration documentation

Supported configurations per body

| Setting        | Allowed values | Required | Description                                         |
| :------------- | :------------- | :------- | :-------------------------------------------------- |
| provider       | string         | yes      | cloudfront, cloudflare                              |
| trustip        | []string       | No       | IP or IP range to trust                             |
| disableDefault | bool           | Yes      | Disable the built in list of CloudFlare IPs/Servers |


### Enable the plugin

```yaml
experimental:
  plugins:
    traefikdisolver:
      modulename: github.com/kyaxcorp/traefikdisolver
      version: v1.0.5
```

### Plugin configuration

```yaml
http:
  middlewares:
    traefikdisolver-cloudfront:
      plugin:
        traefikdisolver:
          provider: cloudfront # cloudfront, cloudflare
          disableDefault: true
          trustip: # Trust IPS not required if disableDefault is false - we will allocate Cloud Flare IPs automatically
            - "0.0.0.0/0"
            - "::/0"
    traefikdisolver-cloudflare:
      plugin:
        traefikdisolver:
          provider: cloudflare # cloudfront, cloudflare
          disableDefault: true
          trustip: # Trust IPS not required if disableDefault is false - we will allocate Cloud Flare IPs automatically
            - "0.0.0.0/0"
            - "::/0"

  routers:
    my-router-cloudfront:
      rule: Path(`/whoami`)
      service: service-whoami
      entryPoints:
        - http
      middlewares:
        - traefikdisolver-cloudfront
    my-router-cloudflare:
      rule: Path(`/whoami`)
      service: service-whoami
      entryPoints:
        - http
      middlewares:
        - traefikdisolver-cloudflare

  services:
    service-whoami:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:5000
```