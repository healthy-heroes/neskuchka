from fastapi import FastAPI

app = FastAPI()


@app.get("/api/ping")
async def root():
    return {"message": "Pong"}


exercises = [
    {"id": 1, "name": "Push-ups"},
    {"id": 2, "name": "Pull-ups"},
    {"id": 3, "name": "Squats"},
]

@app.get("/api/exersices")
async def get_exersices():
    return exercises


@app.get("/api/exersices/{id}")
async def get_exersice(id: int):
    exercise = next((e for e in exercises if e["id"] == id), None)
    return {"data": exercise}


