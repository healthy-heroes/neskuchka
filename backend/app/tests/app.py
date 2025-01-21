from fastapi.testclient import TestClient
from app.main import app
from app.config import get_settings
from app.tests.config import TestSettings


# override settings for tests
def get_settings_override():
    return TestSettings()


app.dependency_overrides[get_settings] = get_settings_override


def get_client():
    return TestClient(app)
