from fastapi import APIRouter
from fastapi.responses import JSONResponse

router = APIRouter()


@router.get("/api/v1/healthcheck", name="Healthcheck", status_code=200)
async def healthcheck():
    return JSONResponse(status_code=200, content={"status": "available"})
