import path from "node:path";
import { fileURLToPath } from "node:url";

import { mergeConfig, defineConfig } from "vitest/config";
import { storybookTest } from "@storybook/experimental-addon-test/vitest-plugin";

import viteConfig from "./vite.config";

const dirname =
  typeof __dirname !== "undefined"
    ? __dirname
    : path.dirname(fileURLToPath(import.meta.url));

export default mergeConfig(
  viteConfig,
  defineConfig({
    test: {
      pool: "threads",

      coverage: {
        provider: "v8",
        reporter: ["text", "json", "html", "lcov"],
        reportsDirectory: "./coverage",
      },

      workspace: [
        {
          extends: true,

          plugins: [
            // The plugin will run tests for the stories defined in your Storybook config
            // See options at: https://storybook.js.org/docs/writing-tests/test-addon#storybooktest
            storybookTest({ configDir: path.join(dirname, ".storybook") }),
          ],

          test: {
            name: "stories",

            browser: {
              enabled: true,
              provider: "playwright",
              instances: [{ browser: "chromium" }],
              headless: true,
            },

            setupFiles: [".storybook/vitest.setup.ts"],
          },
        },
        {
          extends: true,

          test: {
            environment: "node",
            name: "unit",
            include: ["src/**/*.test.ts"],
          },
        },
      ],
    },
  }),
);
