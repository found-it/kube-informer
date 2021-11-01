# Kube-Informer

Build it

```sh
make mac-binary
```

Run it
```sh
./bin/inform
```

All it does right now is simply print out the events as they happen. You should see a lot of `add` events as the cache is built and then if you run

```sh
kubectl run redis --image=redis:4
kubectl delete po redis

kubectl run redis --image=redis:failme

# wait for a bit...
kubectl delete po redis
```

You'll see some add, update, delete events in series.
