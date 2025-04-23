import eslint from "@eslint/js";
import tseslint from "typescript-eslint";

import reactPlugin from "eslint-plugin-react";
import reactHooksPlugin from "eslint-plugin-react-hooks";
import reactRefreshPlugin from "eslint-plugin-react-refresh";
import prettierPlugin from "eslint-plugin-prettier";

export default tseslint.config(
  { ignores: ["dist", "storybook-static", "coverage"] },

  eslint.configs.recommended,
  tseslint.configs.recommended,

  {
    settings: { react: { version: "19" } },

    files: ["**/*.{ts,tsx}"],

    plugins: {
      react: reactPlugin,
      "react-hooks": reactHooksPlugin,
      "react-refresh": reactRefreshPlugin,
      prettier: prettierPlugin,
    },

    rules: {
      ...reactPlugin.configs.recommended.rules,
      ...reactPlugin.configs["jsx-runtime"].rules,

      ...reactHooksPlugin.configs.recommended.rules,
      ...reactRefreshPlugin.configs.recommended.rules,

      ...prettierPlugin.configs.recommended.rules,
    },
  },
);
