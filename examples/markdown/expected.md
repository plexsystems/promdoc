# Alerts

## Rule Groups

* [Description](#description)
* [Message](#message)
* [MultiAlert1](#multialert1)
* [MultiAlert2](#multialert2)

## Description

|Name|Summary|Description|Severity|Expr|For|Runbook|
|---|---|---|---|---|---|---|
|DescriptionAlert|TestSummary|TestDescription|TestSeverity|up == 0|1w|[TestRunbookURL](TestRunbookURL)|

## Message

|Name|Summary|Description|Severity|Expr|For|Runbook|
|---|---|---|---|---|---|---|
|MessageAlert|TestSummary|TestMessage|TestSeverity|api_http_request_latencies_second{quantile="0.5"} > 1|15m|[TestRunbookURL](TestRunbookURL)|

## MultiAlert1

|Name|Summary|Description|Severity|Expr|For|Runbook|
|---|---|---|---|---|---|---|
|Alert1|TestSummary1|TestAlert1|TestSeverity1|job:request_latency_seconds:mean5m{job="myjob"} > 0.5|2d|[TestRunbookURL1](TestRunbookURL1)|

## MultiAlert2

|Name|Summary|Description|Severity|Expr|For|Runbook|
|---|---|---|---|---|---|---|
|Alert2|TestSummary2|TestAlert2|TestSeverity2|(   predict_linear(prometheus_notifications_queue_length{job="prometheus"}[5m], 60 * 30) >   min_over_time(prometheus_notifications_queue_capacity{job="prometheus"}[5m]) )||[TestRunbookURL2](TestRunbookURL2)|
