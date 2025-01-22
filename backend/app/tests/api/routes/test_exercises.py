from fastapi.testclient import TestClient


def test_create_exercise(client: TestClient):
    exercise_data = {
        "name": "Push-ups",
        "slug": "push-ups",
        "description": "Basic push-ups exercise",
    }

    response = client.post("/api/v1/exercises/", json=exercise_data)

    assert response.status_code == 201
    created_exercise = response.json()
    assert created_exercise["name"] == exercise_data["name"]
    assert created_exercise["slug"] == exercise_data["slug"]
    assert created_exercise["description"] == exercise_data["description"]


def test_get_exercises(client: TestClient):
    response = client.get("/api/v1/exercises/")

    assert response.status_code == 200
    exercises = response.json()
    assert isinstance(exercises, list)


def test_get_exercise_by_slug(client: TestClient):
    # First create an exercise
    exercise_data = {
        "name": "Squats",
        "slug": "squats",
        "description": "Basic squats exercise",
    }
    client.post("/api/v1/exercises/", json=exercise_data)

    # Then retrieve it by slug
    response = client.get(f"/api/v1/exercises/{exercise_data['slug']}")

    assert response.status_code == 200
    exercise = response.json()
    assert exercise["name"] == exercise_data["name"]
    assert exercise["slug"] == exercise_data["slug"]
    assert exercise["description"] == exercise_data["description"]


# def test_get_exercise_nonexistent_slug(client: TestClient):
#     response = client.get("/api/v1/exercises/nonexistent_exercise")
#     assert response.status_code == 404


def test_create_exercise_duplicate_slug(client: TestClient):
    # Create first exercise
    exercise_data = {
        "name": "Push-ups",
        "slug": "push-ups_uniq",
        "description": "Basic push-ups exercise",
    }
    response = client.post("/api/v1/exercises/", json=exercise_data)
    assert response.status_code == 201

    # Try to create second exercise with same slug
    duplicate_exercise = {
        "name": "Different Push-ups",
        "slug": "push-ups_uniq",  # Same slug as first exercise
        "description": "Another push-ups variation",
    }
    response = client.post("/api/v1/exercises/", json=duplicate_exercise)
    assert response.status_code == 400
