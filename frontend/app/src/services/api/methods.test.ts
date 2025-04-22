import { describe, expect, it, beforeAll, afterAll, afterEach } from "vitest";
import { http, HttpResponse } from "msw";
import { setupServer } from "msw/node";

import { getMainTrackWorkouts } from "./methods";
import { HttpError } from "./httpErrors";

const endpoints = {
  getMainTrackWorkouts: {
    url: "http://localhost:8080/api/v1/tracks/main/last_workouts",
    data: {
      Workouts: [],
      Exercises: {},
    },
  },
};

describe("getMainTrackWorkouts", () => {
  const server = setupServer();
  beforeAll(() => server.listen());
  afterEach(() => server.resetHandlers());
  afterAll(() => server.close());

  it("should return track workouts when the API call is successful", async () => {
    server.use(
      http.get(endpoints.getMainTrackWorkouts.url, () => {
        return HttpResponse.json(endpoints.getMainTrackWorkouts.data);
      }),
    );

    const result = await getMainTrackWorkouts();
    expect(result).toEqual(endpoints.getMainTrackWorkouts.data);
  });

  it("should throw HttpError when the API call fails", async () => {
    server.use(
      http.get(endpoints.getMainTrackWorkouts.url, () => {
        return new HttpResponse(
          JSON.stringify({ details: "Internal server error" }),
          { status: 500 },
        );
      }),
    );

    await expect(getMainTrackWorkouts()).rejects.toThrow(HttpError);
  });
});
