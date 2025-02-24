const API_URL = "http://localhost:8000/api/v1";

export class HttpError extends Error {
  constructor(
    public status: number,
    public details: object,
  ) {
    super(JSON.stringify(details));
  }
}

export interface Exercise {
  slug: string;
  name: string;
  description: string;
}

export interface Workout {
  id: number;
  date: string;

  sections: Array<{
    title: string;
    protocol: {
      title: string;
      description: string;
    };
    exercises: Array<{
      exercise_slug: string;
    }>;
  }>;
}

export interface TrackWorkouts {
  workouts: Array<Workout>;
  exercises: Record<string, Exercise>;
}

export async function getMainTrackWorkouts(): Promise<TrackWorkouts> {
  const response = await fetch(`${API_URL}/tracks/main/last_workouts`);
  if (!response.ok) {
    throw new HttpError(response.status, await response.json());
  }

  return response.json();
}
