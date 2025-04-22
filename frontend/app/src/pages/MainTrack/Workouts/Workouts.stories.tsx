import { Meta, StoryObj } from "@storybook/react";

import { Workouts } from "./";

export const WorkoutsData = {
  trackWorkouts: {
    Workouts: [
      {
        ID: 1,
        Date: "2025-02-03T00:00:00Z",
        Sections: [
          {
            Title: "Разминка",
            Protocol: {
              Title: "3 раунда",
              Description: "",
            },
            Exercises: [
              {
                ExerciseSlug: "table",
              },
              {
                ExerciseSlug: "forward-bend",
              },
              {
                ExerciseSlug: "plank-with-jumping-jack",
              },
            ],
          },
          {
            Title: "Комплекс",
            Protocol: {
              Title: "5 раундов",
              Description: "",
            },
            Exercises: [
              {
                ExerciseSlug: "deadlift-on-one-leg",
              },
              {
                ExerciseSlug: "squats",
              },
              {
                ExerciseSlug: "push-up-with-hands-over-head",
              },
            ],
          },
        ],
      },
      {
        ID: 2,
        Date: "2025-01-31T00:00:00Z",
        Sections: [
          {
            Title: "Разминка",
            Protocol: {
              Title: "3 раунда",
              Description: "",
            },
            Exercises: [
              {
                ExerciseSlug: "non-exists",
              },
            ],
          },
        ],
      },
    ],
    Exercises: {
      "deadlift-on-one-leg": {
        Slug: "deadlift-on-one-leg",
        Name: "Становая на одной",
        Description: "",
      },
      "forward-bend": {
        Slug: "forward-bend",
        Name: "Наклоны вперед",
        Description: "",
      },
      "plank-with-jumping-jack": {
        Slug: "plank-with-jumping-jack",
        Name: "Качающиеся планки",
        Description: "",
      },
      "push-up-with-hands-over-head": {
        Slug: "push-up-with-hands-over-head",
        Name: "С груди над головой",
        Description: "",
      },
      squats: {
        Slug: "squats",
        Name: "Приседания",
        Description: "",
      },
      table: {
        Slug: "table",
        Name: "Стол",
        Description: "",
      },
    },
  },
};

const meta = {
  component: Workouts,
  title: "Workouts",
  excludeStories: /.*Data$/,
  tags: ["autodocs"],
} satisfies Meta<typeof Workouts>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {
    ...WorkoutsData,
  },
};

export const Empty: Story = {
  args: {
    trackWorkouts: {
      Workouts: [],
      Exercises: {},
    },
  },
};
