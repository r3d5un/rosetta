import logging

import structlog
from fastapi import FastAPI, Request
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor

from src.python.api.healtcheck import router as healthcheck_router
from src.python.core.config import read_config

app = FastAPI()
cfg = read_config()

for name in logging.root.manager.loggerDict:
    logging.getLogger(name).setLevel(logging.CRITICAL)
    logging.getLogger(name).handlers = []

# Configure structlog to intercept standard logging
structlog.configure(
    processors=[
        structlog.stdlib.filter_by_level,
        structlog.processors.TimeStamper(fmt="iso"),
        structlog.processors.JSONRenderer(),
    ],
    context_class=dict,
    logger_factory=structlog.stdlib.LoggerFactory(),
    wrapper_class=structlog.stdlib.BoundLogger,
    cache_logger_on_first_use=True,
)

# Wrap the standard logging configuration with structlog
logging.basicConfig(
    format="%(message)s",
    level=logging.INFO,
)

# Replace the standard logging handlers with structlog handlers
logging.getLogger().handlers = [logging.StreamHandler()]

# Get a structlog logger
logger = structlog.get_logger()


@app.get("/")
async def root():
    return {"message": "Hello, World!"}


@app.middleware("http")
async def log_request_middleware(request: Request, call_next):
    logger.info("request received", method=request.method, url=request.url)
    response = await call_next(request)
    logger.info("request completed", method=request.method, url=request.url)
    return response


app.include_router(healthcheck_router)
FastAPIInstrumentor.instrument_app(app)
