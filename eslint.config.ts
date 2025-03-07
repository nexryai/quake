import { defineConfigWithVueTs, vueTsConfigs } from "@vue/eslint-config-typescript";
// @ts-ignore
import pluginImport from "eslint-plugin-import";
import pluginVue from "eslint-plugin-vue";
import pluginTypescript from "typescript-eslint";

// To allow more languages other than `ts` in `.vue` files, uncomment the following lines:
// import { configureVueProject } from '@vue/eslint-config-typescript'
// configureVueProject({ scriptLangs: ['ts', 'tsx'] })
// More info at https://github.com/vuejs/eslint-config-typescript/#advanced-setup

export default defineConfigWithVueTs(
    {
        name: "app/files-to-lint",
        files: ["**/*.{ts,mts,tsx,vue}"],
    },

    {
        name: "app/files-to-ignore",
        ignores: ["**/dist/**", "**/dist-ssr/**", "**/coverage/**"],
    },

    pluginVue.configs["flat/essential"],
    vueTsConfigs.recommended,
    pluginImport.flatConfigs.recommended,
    pluginImport.flatConfigs.typescript,

    {
        files: ["**/*.vue", "**/*.ts", "**/*.js"],

        languageOptions: {
            parserOptions: {
                parser: pluginTypescript.parser
            }
        },
        rules: {
            "@typescript-eslint/no-unused-expressions": "off",
            "@typescript-eslint/ban-ts-comment": "off",
            "@typescript-eslint/no-explicit-any": "off",
            "indent": ["error", 4, { "SwitchCase": 1 }],
            "quotes": ["error", "double"],
            "semi": ["error", "always"],
            "import/no-cycle": "error",
            // https://github.com/import-js/eslint-plugin-import/issues/2765
            "import/no-unresolved": "off",
            "import/no-named-as-default": "off",
            "import/order": [
                "error",
                {
                    groups: ["builtin", "external", "parent", "sibling", "index", "object"],
                    pathGroups: [
                        {
                            pattern: "{svelte,$app/**,elysia}",
                            group: "builtin",
                            position: "before",
                        },
                        {
                            pattern: "{@src/**,$lib/**,@prisma/client}",
                            group: "parent",
                            position: "before",
                        },
                    ],
                    pathGroupsExcludedImportTypes: ["type"],
                    alphabetize: {
                        order: "asc",
                    },
                    "newlines-between": "always",
                },
            ],
        }
    },
);
