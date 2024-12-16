import { Stack } from "expo-router";
import * as SplashScreen from "expo-splash-screen";
import { PaperProvider } from "react-native-paper";
import { useEffect } from "react";

// Initialize localization.
import "@/components/i18n";

// Prevent the splash screen from auto-hiding before asset loading is complete.
SplashScreen.preventAutoHideAsync();

export default function RootLayout() {
  useEffect(() => {
    SplashScreen.hideAsync();
  }, []);

  return (
    <PaperProvider>
      <Stack>
        <Stack.Screen name="(tabs)" options={{ headerShown: false }} />
      </Stack>
    </PaperProvider>
  );
}
