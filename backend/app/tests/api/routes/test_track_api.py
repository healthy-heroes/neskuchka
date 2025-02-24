import pytest
from app.domain.exercise import ExerciseRepository
from app.domain.track import Track, TrackRepository
from app.domain.workout import Workout, WorkoutRepository
from app.tests.fixtures.domain import create_track, create_workout


def track_for_assert(track: Track) -> dict:
    track_assert = track.model_dump()
    track_assert["id"] = str(track.id)
    track_assert["owner_id"] = str(track.owner_id)
    return track_assert


def workout_for_assert(workout: Workout) -> dict:
    workout_assert = workout.model_dump()
    workout_assert["id"] = str(workout.id)
    workout_assert["date"] = str(workout.date)
    workout_assert["track_id"] = str(workout.track_id)
    return workout_assert


def test_get_main_track(client, track_repo: TrackRepository):
    track = create_track()
    track_repo.add(track)

    response = client.get("/api/v1/tracks/main")

    assert response.status_code == 200
    assert response.json() == track_for_assert(track)


def test_get_non_existent_main_track(client, track_repo: TrackRepository):
    response = client.get("/api/v1/tracks/main")

    assert response.status_code == 404


@pytest.mark.skip(reason="TODO: пока сломано, нужно расширить фикстуры")
def test_get_workouts_for_main_track(
    client, track_repo: TrackRepository, workout_repo: WorkoutRepository
):
    track = create_track()
    track_repo.add(track)

    workouts = []
    for _ in range(12):
        workout = create_workout(track_id=track.id)
        workout_repo.add(workout)
        workouts.append(workout)

    response = client.get("/api/v1/tracks/main/last_workouts")

    assert response.status_code == 200
    # todo: пока захардкожен лимит
    assert response.json() == [workout_for_assert(workout) for workout in workouts[:10]]


def test_get_workouts_for_non_existent_track(
    client,
    track_repo: TrackRepository,
    workout_repo: WorkoutRepository,
    exercise_repo: ExerciseRepository,
):
    response = client.get("/api/v1/tracks/main/last_workouts")

    assert response.status_code == 404
