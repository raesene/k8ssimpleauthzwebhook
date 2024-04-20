# K8s Simple Authz Webhook

This is a very simplistic implementation of a Kubernetes authorization webhook, for test purposes. Do Not use in production!

The structure is that there's a file called `rights.txt` which contains rights assigned to users one per line. This is read by the code which then listens on port `8888/TCP` for incoming request.



## Testing the API


Testing a valid request

```
curl -X POST -H "Content-Type: application/json" -d @validtestrequest.json http://localhost:8888/authorize
```

Testing an invalid request

```
curl -X POST -H "Content-Type: application/json" -d @invalidtestrequest.json http://localhost:8888/authorize
```


## References

[Kubernetes documentation on Authorization Webhooks](https://kubernetes.io/docs/reference/access-authn-authz/webhook/)
[k8s webhook helloworld](https://github.com/salrashid123/k8s_webhook_helloworld)
[auth webhook sample](https://github.com/dinumathai/auth-webhook-sample)

