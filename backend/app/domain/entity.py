from pydantic import BaseModel, ConfigDict


class EntityModel(BaseModel):
    """
    Базовая доменная модель

    Все доменные модели будут потомками, будут наследовать базовые свойства, например, неизменяемость
    """
    model_config = ConfigDict(frozen=True)


# todo: Какой uuid использовать?