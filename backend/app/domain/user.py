from uuid import UUID
from typing import Annotated
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


class UserRepository(EntityModel):
    def get_by_id(self, user_id: UserId) -> User | None:
        pass
