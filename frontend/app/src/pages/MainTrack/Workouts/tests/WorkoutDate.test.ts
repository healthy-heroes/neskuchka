import { describe, it, expect } from "vitest";

import { WorkoutDate } from "../WorkoutDate";

describe("WorkoutDate", () => {
  it("should return the correct date without year if it is current year", () => {
    const date = `${new Date().getFullYear()}-03-04`;
    const formattedDate = WorkoutDate(date);
    expect(formattedDate).toBe("4 марта");
  });

  it("should return the correct date with year", () => {
    const date = `${new Date().getFullYear() - 2}-02-28`;
    const formattedDate = WorkoutDate(date);
    expect(formattedDate).toBe("28 февраля 2023");
  });
});
