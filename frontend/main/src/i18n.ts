import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import LanguageDetector from "i18next-browser-languagedetector";

import enJSONLogin from "./assets/locales/en/login.json";
import ruJSONLogin from "./assets/locales/ru/login.json";

import enJSONConfirm from "./assets/locales/en/confirm.json";
import ruJSONConfirm from "./assets/locales/ru/confirm.json";

await i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    debug: true,
    fallbackLng: "en",
    interpolation: {
      escapeValue: false,
    },
    resources: {
      en: { ...enJSONConfirm, ...enJSONLogin },
      ru: { ...ruJSONConfirm, ...ruJSONLogin },
    },
  });
