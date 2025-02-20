from contextlib import asynccontextmanager
from fastapi import FastAPI
from pydantic import BaseModel
from sqlmodel import SQLModel

from app.api.main import init_api
from app.config.main import settings
from app.infrastructure.db.database import db


@asynccontextmanager
async def lifespan(app: FastAPI):
    print("app:Starting up")
    SQLModel.metadata.create_all(db.engine)

    yield

    print("app:Shutting down")
    db.dispose()


app = FastAPI(lifespan=lifespan)


# attach routers
init_api(app, settings)


# root endpoint
class Info(BaseModel):
    name: str
    environment: str


@app.get("/")
def info():
    return Info(name=settings.app_name, environment=settings.environment)
