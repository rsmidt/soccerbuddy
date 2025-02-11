export default {
  expo: {
    name: "soccerbuddy",
    slug: "soccerbuddy",
    version: "1.0.0",
    orientation: "portrait",
    icon: "./assets/images/icon.png",
    scheme: "soccerbuddy",
    userInterfaceStyle: "automatic",
    newArchEnabled: true,
    ios: {
      supportsTablet: false,
      entitlements: {
        "aps-environment": "production",
      },
      googleServicesFile:
        process.env.GOOGLE_SERVICES_PLIST ?? "./GoogleService-Info.plist",
      bundleIdentifier: "dev.rsmidt.soccerbuddy",
    },
    android: {
      adaptiveIcon: {
        foregroundImage: "./assets/images/adaptive-icon.png",
        backgroundColor: "#1AEB28",
      },
      package: "dev.rsmidt.soccerbuddy",
      googleServicesFile:
        process.env.GOOGLE_SERVICES_JSON ?? "./google-services.json",
    },
    web: {
      bundler: "metro",
      output: "static",
      favicon: "./assets/images/favicon.png",
    },
    notification: {
      icon: "./assets/images/notification-icon.png",
    },
    plugins: [
      "expo-router",
      [
        "expo-splash-screen",
        {
          image: "./assets/images/splash-icon.png",
          imageWidth: 200,
          resizeMode: "contain",
          backgroundColor: "#ffffff",
        },
      ],
      "expo-localization",
      "expo-secure-store",
      "@react-native-firebase/app",
      "@react-native-firebase/messaging",
      [
        "expo-build-properties",
        {
          ios: {
            useFrameworks: "static",
          },
          android: {
            extraMavenRepos: [
              "../../node_modules/@notifee/react-native/android/libs",
            ],
          },
        },
      ],
    ],
    experiments: {
      typedRoutes: true,
    },
    updates: {
      enabled: false,
      url: "https://u.expo.dev/8227016f-4134-48e1-8798-4343cf0c04e1",
    },
    extra: {
      router: {
        origin: false,
      },
      eas: {
        projectId: "8227016f-4134-48e1-8798-4343cf0c04e1",
      },
    },
    runtimeVersion: {
      policy: "appVersion",
    },
  },
};
