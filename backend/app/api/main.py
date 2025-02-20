from fastapi import APIRouter, FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.config.default import BaseSettings
from app.api.routes import utils, exercises, tracks

api_router = APIRouter()

api_router.include_router(exercises.router)
api_router.include_router(tracks.router)
api_router.include_router(utils.router)


def init_api(app: FastAPI, settings: BaseSettings):
    app.add_middleware(
        CORSMiddleware,
        allow_origins=["*"],
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )

    # attach routers
    app.include_router(api_router, prefix=settings.API_V1_PREFIX)
