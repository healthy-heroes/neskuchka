from sqlmodel import Session

from app.domain.user import User
from app.domain.track import Track
from app.infrastructure.db.track import TrackDbRepository


def test_add_get_track(session: Session):
    user_id = User.create_id("test-user")
    track = Track(
        name="Test track",
        owner_id=user_id,
    )

    track_repository = TrackDbRepository(session)

    track_from_db = track_repository.get_main_track()
    assert track_from_db is None, "Should be no main track if not added"

    track_repository.add(track)

    track_from_db = track_repository.get_by_id(track.id)
    assert track_from_db == track, "Should be able to get track by id"

    track_from_db = track_repository.get_main_track()
    assert track_from_db == track, "Should be able to get main track"
