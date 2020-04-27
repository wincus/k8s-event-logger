# k8s-event-logger

## How to use

Apply the files under `/kubernetes` and then run:

```bash
kubectl logs -l name=k8s-event-logger-pod -f
```