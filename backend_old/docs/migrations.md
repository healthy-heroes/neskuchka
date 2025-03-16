# Миграции

## Концепция миграции домена

Миграции должны жить в соответствующем инфраструктурном модуле

Данные можно мигрировать путем добавления в таблицу поля с версией.
По версии можно проводить синрохронную и асинхронную миграцию.

Синхронная миграция - это написание миграционного файла, в котором обходим данные и обновляем их. Похожа на асинхронный пример (ниже), только проходятся по всем данным.

Асинхронная миграция - это миграции на основе версии, которая происходит в момент получения данных из базы. Пример:

```
CURRENT_VERSION = 1

class WorkoutModel(SQLModel, table=True):
    __tablename__ = "workout"

    id: str = Field(default=None, primary_key=True)
    sections: str
    schema_version: int = Field(default=1)  # Add version field

    def to_domain(self) -> Workout:
        sections_data = loads(self.sections)
        
        # Migrate data structure if needed
        sections = self._migrate_sections(sections_data, self.schema_version)
        
        return Workout(
            id=WorkoutId(self.id),
            sections=sections["sections"]
        )

# 

class WorkoutDbRepository(WorkoutRepository):
    def __init__(self, session: Session):
        self.session = session

    def add(self, workout: Workout) -> Workout:
        db_workout = WorkoutModel(
            id=str(workout.id),
            sections=workout.model_dump_json(include={'sections'}),
            schema_version=CURRENT_VERSION  # Current version
        )
```

### Сценарии для миграции

Добавление новых полей. Можно делать любым способом, так как новые поля будут просто игнорироваться.

Удаление поля. Можно безопасно удалить поле. После выкатки домена прокатить миграцию.

Изменение текущего поля. Нужно будет сначала выкатить новый домен и одновременно поддержать асинхронную миграцию. После выкатки домена можно прокатить синхронную миграцию.
