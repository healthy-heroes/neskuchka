import { Meta, StoryObj } from "@storybook/react";

import { LandingPage } from ".";

const meta: Meta<typeof LandingPage> = {
  component: LandingPage,
};

export default meta;

type Story = StoryObj<typeof LandingPage>;

export const Default: Story = {};
