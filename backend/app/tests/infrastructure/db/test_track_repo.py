import uuid
from sqlmodel import Session

from app.domain.user import UserId
from app.domain.track import Track, TrackId
from app.infrastructure.db.track import TrackDbRepository


def test_add_get_track(session: Session):
    user_id = UserId(str(uuid.uuid5(uuid.NAMESPACE_DNS, "test-user")))
    track_id = TrackId(str(uuid.uuid5(uuid.NAMESPACE_DNS, "test-track")))
    track = Track(
        id=track_id,
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
