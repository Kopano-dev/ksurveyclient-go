# Kopano Survey Client for Go

This module implements a client to submit survey data via HTTPS to a remote 
service.

By default, data is transmitted to the stats service operated by Kopano at
`https://stats.kopano.io/api/stats/v1/submit`. Stats are transmitted upon
initialization and afterwards in a one hour interval.



