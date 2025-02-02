import uuid
from sqlmodel import Session

from app.domain.user import User, UserId
from app.infrastructure.db.user import UserDbRepository


def test_add_get_user(session: Session):
    user_id = UserId(str(uuid.uuid5(uuid.NAMESPACE_DNS, "test-user")))

    user = User(
        id=user_id,
        name="Test User",
        login="test_user",
        email="test_user@example.com",
    )

    user_repository = UserDbRepository(session)
    user_repository.add(user)

    user_from_db = user_repository.get_by_id(user_id)
    assert user_from_db == user
