import uuid
from fastapi.testclient import TestClient
from sqlmodel import Session

from app.infrastructure.db.track import TrackModel
from app.infrastructure.db.user import UserModel


def test_get_main_track(client: TestClient, session: Session):
    user_id = str(uuid.uuid5(uuid.NAMESPACE_DNS, "test-user"))
    db_user = UserModel(
        id=user_id,
        name="Test User",
        login="test_user",
        email="test_user@example.com",
    )
    session.add(db_user)
    session.commit()

    track_id = str(uuid.uuid5(uuid.NAMESPACE_DNS, "main-track"))
    db_track = TrackModel(
        id=track_id,
        name="Main Track",
        owner_id=user_id,
    )
    session.add(db_track)
    session.commit()

    response = client.get("/api/v1/tracks/main")
    assert response.status_code == 200
    assert response.json() == {
        "id": track_id,
        "name": "Main Track",
        "owner_id": user_id,
    }
