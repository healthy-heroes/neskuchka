import React from "react";
import type { Preview } from "@storybook/react";
import {
  Provider as SpectrumProvider,
  defaultTheme,
} from "@adobe/react-spectrum";

const preview: Preview = {
  parameters: {
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/i,
      },
    },
  },

  decorators: [
    (Story) => (
      <SpectrumProvider theme={defaultTheme} colorScheme="light" scale="large">
        <Story />
      </SpectrumProvider>
    ),
  ],
};

export default preview;
