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

### Notes re CloudFlare

One thing included in this plugin is we bundle the CloudFlare server IPs with it, so you do not have to define them manually.  
However on the flip-side, if you want to, you can just disable them by setting `disableDefault` to `true`.

If you do not define `trustip` and `disableDefault`, it doesn't seem to load the plugin, so just set `disableDefault` to `false` and you are able to use the default IP list.

### Enable the plugin

```yaml
experimental:
  plugins:
    traefik-client-real-ip:
      modulename: github.com/kyaxcorp/traefik-client-real-ip
      version: v1.0.0
```

### Plugin configuration

```yaml
http:
  middlewares:
    client-real-ip:
      plugin:
        traefik-client-real-ip:
          provider: cloudfront # cloudfront, cloudflare
          disableDefault: true
          trustip: # Trust IPS not required if disableDefault is false - we will allocate Cloud Flare IPs automatically
            - "0.0.0.0/0"
            - "::/0"

  routers:
    my-router:
      rule: Path(`/whoami`)
      service: service-whoami
      entryPoints:
        - http
      middlewares:
        - client-real-ip

  services:
    service-whoami:
      loadBalancer:
        servers:
          - url: http://127.0.0.1:5000
```