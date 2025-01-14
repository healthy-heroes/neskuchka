from fastapi.testclient import TestClient

from .config import get_settings, TestSettings
from .main import app

client = TestClient(app)


def get_settings_override():
    return TestSettings()


app.dependency_overrides[get_settings] = get_settings_override


def test_info():
    response = client.get("/")
    assert response.status_code == 200
    assert response.json() == {"name": "Test app", "environment": "test"}
