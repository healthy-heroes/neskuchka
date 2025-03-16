from sqlmodel import Session

from app.tests.fixtures.domain import create_user
from app.infrastructure.db.user import UserDbRepository


def test_add_get_user(session: Session):
    user = create_user()

    user_repository = UserDbRepository(session)
    user_repository.add(user)

    user_from_db = user_repository.get_by_id(user.id)
    assert user_from_db == user
