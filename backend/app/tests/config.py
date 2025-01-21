from app.config import Settings


class TestSettings(Settings):
    app_name: str = "Test app"

    environment: str = "test"

    API_V1_PREFIX: str = "/api/v1"


test_settings = TestSettings()
