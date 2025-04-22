import { describe, it, expect } from "vitest";

import { formatIsoDate } from "./dates";

describe("Utils: dates.formatIsoDate", () => {
  it("should return the correct date without year if it is current year", () => {
    const date = `${new Date().getFullYear()}-03-04`;
    const formattedDate = formatIsoDate(date);
    expect(formattedDate).toBe("4 марта");
  });

  it("should return the correct date with year", () => {
    const date = `${new Date().getFullYear() - 2}-02-28`;
    const formattedDate = formatIsoDate(date);
    expect(formattedDate).toBe("28 февраля 2023");
  });
});
