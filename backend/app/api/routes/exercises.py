from fastapi import APIRouter
from pydantic import BaseModel

router = APIRouter(prefix="/exercises", tags=["exercises"])


class Exercise(BaseModel):
    id: int
    name: str
    description: str | None = None


exercises = [
    Exercise(id=1, name="Push-ups"),
    Exercise(id=2, name="Pull-ups"),
    Exercise(id=3, name="Squats"),
    Exercise(id=4, name="SvinKy"),
]


@router.post("/")
async def create_exercise(item: Exercise):
    exercises.append(item)
    return item


@router.get("/")
async def get_exersices(skip: int = 0):
    return exercises[skip:]


@router.get("/{id}")
async def get_exersice(id: int):
    exercise = next((e for e in exercises if e["id"] == id), None)
    return {"data": exercise}
