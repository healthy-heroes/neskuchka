import { fn } from "@storybook/test";
import * as actual from "./methods";

export * from "./methods";
export const getMainTrackWorkouts = fn(actual.getMainTrackWorkouts).mockName(
  "getMainTrackWorkouts",
);
