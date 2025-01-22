from app.config.default import AppSettings, DatabaseEngineArgs, DatabaseSettings


class TestAppSettings(AppSettings):
    app_name: str = "Test app"
    environment: str = "test"

    database: DatabaseSettings = DatabaseSettings(
        dsn="sqlite:///database_test.db",
        engine_args=DatabaseEngineArgs(
            connect_args={"check_same_thread": False},
        ),
    )
