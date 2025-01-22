from fastapi.testclient import TestClient


def test_info(client: TestClient):
    response = client.get("/")
    assert response.status_code == 200
    assert response.json() == {"name": "Test app", "environment": "test"}
