# Rosetta Python

Running the application with telemetry:

```bash
export OTEL_PYTHON_LOGGING_AUTO_INSTRUMENTATION_ENABLED=true
poetry run opentelemetry-instrument --traces_exporter console --metrics_exporter console --logs_exporter console --service_name rosetta fastapi dev src/python/main.py
```
