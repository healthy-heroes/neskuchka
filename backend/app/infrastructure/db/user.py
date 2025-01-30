from sqlmodel import Field, SQLModel, Session, select

from app.domain.user import UserRepository, User, UserId


class UserModel(SQLModel, table=True):
    __tablename__ = "user"

    id: str = Field(default=None, primary_key=True)
    name: str
    login: str
    email: str

    def to_domain(self) -> User:
        return User(
            id=UserId(self.id),
            name=self.name,
            login=self.login,
            email=self.email,
        )


class UserDbRepository(UserRepository):
    def __init__(self, session: Session):
        self.session = session

    def get_by_id(self, user_id: UserId) -> User | None:
        query = select(UserModel).where(UserModel.id == user_id).limit(1)
        result = self.session.exec(query).first()

        if not result:
            return None

        return result.to_domain()
