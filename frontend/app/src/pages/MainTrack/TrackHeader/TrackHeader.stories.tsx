import { Meta, StoryObj } from "@storybook/react";

import { TrackHeader } from ".";

const meta = {
  component: TrackHeader,
  title: "TrackHeader",
  tags: ["autodocs"],
} satisfies Meta<typeof TrackHeader>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
  args: {},
};
