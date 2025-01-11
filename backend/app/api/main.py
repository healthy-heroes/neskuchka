from fastapi import APIRouter

from app.api.routes import utils, exercises

api_router = APIRouter()

api_router.include_router(exercises.router)
api_router.include_router(utils.router)
