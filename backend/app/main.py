from contextlib import asynccontextmanager
from fastapi import FastAPI
from pydantic import BaseModel

from app.api.main import api_router
from app.infrastructure.db.session import create_db_and_tables
from app.config.main import settings


@asynccontextmanager
async def lifespan(app: FastAPI):
    # Startup: connect to db
    # await db.connect()
    print("app:Starting up")
    create_db_and_tables()

    yield

    # Shutdown: disconnect from db
    # await db.disconnect()
    print("app:Shutting down")


app = FastAPI(lifespan=lifespan)


# attach routers
app.include_router(api_router, prefix=settings.API_V1_PREFIX)


# root endpoint
class Info(BaseModel):
    name: str
    environment: str


@app.get("/")
def info():
    return Info(name=settings.app_name, environment=settings.environment)
