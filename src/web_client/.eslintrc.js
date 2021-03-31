module.exports = {
    root: true,
    parser: "@typescript-eslint/parser",
    parserOptions: {
        ecmaFeatures: {
            jsx: true,
        },
        ecmaVersion: 2018,
        sourceType: "module",
    },
    extends: [
        "eslint:recommended",
        "plugin:react/recommended",
        "plugin:@typescript-eslint/recommended",
        "plugin:@typescript-eslint/eslint-recommended",
        "prettier"
    ],
    env: {
        browser: true,
        es6: true,
    },
    globals: {
        Atomics: "readonly",
        SharedArrayBuffer: "readonly",
    },
    plugins: ["react", "react-hooks", "@typescript-eslint"],
    rules: {
        "@typescript-eslint/explicit-module-boundary-types": 0,
        "@typescript-eslint/no-var-requires": 0,
        "react/prop-types": 0,
    },
    settings: {
        react: {
            version: "detect",
        },
    },
};
