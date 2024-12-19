import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import en from "@/i18n/en/translation.json";
import de from "@/i18n/de/translation.json";
import { getLocales } from "expo-localization";

export const locale = getLocales()[0].languageCode ?? "en";

i18n.use(initReactI18next).init({
  lng: locale,
  fallbackLng: "en",
  interpolation: {
    escapeValue: false,
  },
  resources: {
    en: {
      translation: en,
    },
    de: {
      translation: de,
    },
  },
});

export default i18n;
