from fastapi import FastAPI

app = FastAPI()


@app.get("/api/ping")
async def root():
    return {"message": "Pong"}