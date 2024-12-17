module.exports = {
  input: [
    "!app/**/*.{ts,tsx}",
    "!components/**/*.{ts,tsx}",
    "tmp/**/*.{js,jsx}",
    "!components/**/*.{ts,tsx}",
    "!app/**/*.spec.{ts,tsx}",
    "!**/node_modules/**",
  ],
  output: "./",
  options: {
    removeUnusedKeys: true,
    debug: false,
    func: {
      list: ["i18next.t", "i18n.t", "t"],
      extensions: [".js", ".jsx"],
    },
    trans: {
      component: "Trans",
      i18nKey: "i18nKey",
      defaultsKey: "defaults",
      extensions: [".js", ".jsx"],

      // https://react.i18next.com/latest/trans-component#usage-with-simple-html-elements-like-less-than-br-greater-than-and-others-v10.4.0
      supportBasicHtmlNodes: true, // Enables keeping the name of simple nodes (e.g. <br/>) in translations instead of indexed keys.
      keepBasicHtmlNodesFor: ["br", "strong", "i", "p"], // Which nodes are allowed to be kept in translations during defaultValue generation of <Trans>.
    },
    lngs: ["en", "de"],
    defaultValue: "__STRING_NOT_TRANSLATED__",
    resource: {
      loadPath: "i18n/{{lng}}/{{ns}}.json",
      savePath: "i18n/{{lng}}/{{ns}}.json",
      jsonIndent: 2,
      lineEnding: "\n",
    },
    nsSeparator: false, // namespace separator
    keySeparator: false, // key separator
    interpolation: {
      prefix: "{{",
      suffix: "}}",
    },
    allowDynamicKeys: false,
  },
};
