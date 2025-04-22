import { QueryClientProvider } from "@tanstack/react-query";

import { Meta, StoryObj } from "@storybook/react";

import queryClient from "#api/client";
import { getMainTrackWorkouts } from "#api/methods.mock";

import { MainTrackPage } from ".";
import { WorkoutsData } from "./Workouts/Workouts.stories";
import { HttpError } from "#api/httpErrors";

const meta: Meta<typeof MainTrackPage> = {
  component: MainTrackPage,
  decorators: [
    (Story) => (
      <QueryClientProvider client={queryClient}>
        <Story />
      </QueryClientProvider>
    ),
  ],
};
export default meta;

type Story = StoryObj<typeof MainTrackPage>;

export const Default: Story = {
  beforeEach: async () => {
    getMainTrackWorkouts.mockResolvedValue(WorkoutsData.trackWorkouts);
  },
};

export const Error: Story = {
  beforeEach: async () => {
    getMainTrackWorkouts.mockRejectedValue(
      new HttpError(500, { message: "Internal server error" }),
    );
  },
};

export const Loading: Story = {
  args: {
    isLoading: true,
  },
  beforeEach: async () => {
    getMainTrackWorkouts.mockReturnValue(new Promise(() => {}));
  },
};
