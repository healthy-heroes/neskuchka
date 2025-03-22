const API_URL = "http://localhost:8080/api/v1";

export class HttpError extends Error {
  constructor(
    public status: number,
    public details: object,
  ) {
    super(JSON.stringify(details));
  }
}

export interface Exercise {
  Slug: string;
  Name: string;
  Description: string;
}

export interface Workout {
  ID: number;
  Date: string;

  Sections: Array<{
    Title: string;
    Protocol: {
      Title: string;
      Description: string;
    };
    Exercises: Array<{
      ExerciseSlug: string;
    }>;
  }>;
}

export interface TrackWorkouts {
  Workouts: Array<Workout>;
  Exercises: Record<string, Exercise>;
}

export async function getMainTrackWorkouts(): Promise<TrackWorkouts> {
  const response = await fetch(`${API_URL}/tracks/main/last_workouts`);
  if (!response.ok) {
    throw new HttpError(response.status, await response.json());
  }

  return response.json();
}
