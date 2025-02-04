from sqlmodel import Session

from app.domain.user import User
from app.infrastructure.db.user import UserDbRepository


def test_add_get_user(session: Session):
    user_login = "test-user"
    user_id = User.create_id(user_login)

    user = User(
        id=user_id,
        name="Test User",
        login=user_login,
        email="test_user@example.com",
    )

    user_repository = UserDbRepository(session)
    user_repository.add(user)

    user_from_db = user_repository.get_by_id(user_id)
    assert user_from_db == user
