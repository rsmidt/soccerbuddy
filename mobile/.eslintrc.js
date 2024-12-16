// https://docs.expo.dev/guides/using-eslint/
module.exports = {
  extends: ["expo", "prettier"],
  plugins: ["prettier"],
  ignorePatterns: ["/dist/*", "/app-example"],
  rules: {
    "prettier/prettier": "warn",
  },
};
