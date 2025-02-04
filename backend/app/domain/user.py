from uuid import UUID
from typing import Annotated
import uuid
from pydantic import ConfigDict, EmailStr
from pydantic.types import UuidVersion

from app.domain.entity import EntityModel

UserId = Annotated[UUID, UuidVersion(5)]


class User(EntityModel):
    """
    Модель пользователя сервиса
    """

    model_config = ConfigDict(frozen=True)

    id: UserId
    name: str
    login: str
    email: EmailStr

    @staticmethod
    def create_id(login: str) -> UserId:
        return UserId(str(uuid.uuid5(uuid.NAMESPACE_URL, login)))


class UserRepository:
    def add(self, user: User) -> User:
        pass

    def get_by_id(self, user_id: UserId) -> User | None:
        pass
