from app.tests.app import get_client


client = get_client()


def test_info():
    response = client.get("/")
    assert response.status_code == 200
    assert response.json() == {"name": "Test app", "environment": "test"}
