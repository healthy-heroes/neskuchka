from fastapi import FastAPI
from pydantic import BaseModel

app = FastAPI()


@app.get("/api/ping")
async def root():
    return {"message": "Pong"}


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

@app.post("/api/exercises")
async def create_exercise(item: Exercise):
    exercises.append(item)
    return item

@app.get("/api/exercises")
async def get_exersices(skip: int = 0):
    return exercises[skip:]


@app.get("/api/exercises/{id}")
async def get_exersice(id: int):
    exercise = next((e for e in exercises if e["id"] == id), None)
    return {"data": exercise}


