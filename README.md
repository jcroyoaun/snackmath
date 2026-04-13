# snackmath

Compare calorie density between two foods. Normalize any two items to 100g so you can see which one is the better deal.

PWA -- installable on your phone, works offline (or should). 

Deployed to snackmath.cc

## Dev

```
make run
```

Open http://localhost:7500

Live reload:

```
make run/live
```

## Test

```
make test
```

## Build & push

```
docker buildx build --platform linux/amd64 -t jcroyoaun/snackmath:latest --push .
```

## Deploy (Helm + Istio)

Install:

```
helm install snackmath ./charts/snackmath -n snackmath --create-namespace
```

Upgrade:

```
helm upgrade snackmath ./charts/snackmath -n snackmath
```

Check values in `charts/snackmath/values.yaml` for image tag, domain, and gateway config.

## Config

| Env var | Default | Description |
| --- | --- | --- |
| `HTTP_PORT` | `7500` | Port the server listens on |
| `BASE_URL` | `http://localhost:7500` | Base URL for the app |

---

Running on a Linode Kubernetes cluster because why not.

[github.com/jcroyoaun/snackmath](https://github.com/jcroyoaun/snackmath) · [x.com/jcroyoaun](https://x.com/jcroyoaun)
