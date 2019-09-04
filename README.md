# Kopano Survey Client for Go

This module implements a client to submit survey data via HTTPS to a remote
service.

By default, data is transmitted to the stats service operated by Kopano at
`https://stats.kopano.io/api/stats/v1/submit`. Stats are transmitted upon
initialization and afterwards in a one hour interval.

## Configuration

The survey clients operation is controlled by several environment variables.

```
KOPANO_SURVEYCLIENT_URL
KOPANO_SURVEYCLIENT_START_DELAY
KOPANO_SURVEYCLIENT_ERROR_DELAY
KOPANO_SURVEYCLIENT_INTERVAL
KOPANO_SURVEYCLIENT_INSECURE
KOPANO_SURVEYCLIENT_ENABLED
KOPANO_SURVEYCLIENT_AUTOSURVEY
```

The meaning should be self explaining. To disable all survey operation, set
KOPANO_SURVEYCLIENT_ENABLED to `false` or `no`. To disable the automatic start
of a default survey client, set KOPANO_SURVEYCLIENT_AUTOSURVEY to `false` or
`no`.

## Integration

The API of this model is loosely modelled after [prometheus/client_golang/prometheus](https://github.com/prometheus/client_golang)
and offers various levels of integration. Easiest is to import the `autosurvey`
submodule to enable automatic gathering and submitting of the default collectors
on application start.

```go
import (
	_ "stash.kopano.io/kgol/ksurveyclient-go/autosurvey"
)
```
