from fastapi import FastAPI

from src.python.api.healtcheck import router as healthcheck_router

app = FastAPI()


@app.get("/")
async def root():
    return {"message": "Hello, World!"}


app.include_router(healthcheck_router)
