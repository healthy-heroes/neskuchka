from app.config.default import AppSettings


class TestAppSettings(AppSettings):
    app_name: str = "Test app"

    environment: str = "test"
