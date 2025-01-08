from typing import NewType
from pydantic import EmailStr
from app.domain.entity import EntityModel


# todo: Нужно заменить на uuid
UserId = NewType("UserId", int)


class User(EntityModel):
    """
    Модель пользователя сервиса
    """
    id: UserId
    name: str
    login: str
    email: EmailStr


class RepositoryUser(EntityModel):
    def add(user: User) -> User:
        pass

    def get_by_id(iser_id: UserId) -> User:
        pass