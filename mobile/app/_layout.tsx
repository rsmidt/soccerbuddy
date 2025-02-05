import { Stack } from "expo-router";
import * as SplashScreen from "expo-splash-screen";
import { MD3LightTheme, PaperProvider, useTheme } from "react-native-paper";
import React, { useEffect } from "react";
import Header from "@/components/header";
import { Provider as ReduxProvider } from "react-redux";
import { persistor, store } from "@/store";
import * as SecureStore from "expo-secure-store";
import "@/components/fcm";
// Initialize localization.
import "@/components/i18n";
import {
  validateStoredToken,
  setUnauthenticated,
} from "@/components/auth/auth-slice";
import { SESSION_TOKEN_KEY } from "@/components/auth/constants";
import { StatusBar } from "expo-status-bar";
import { View } from "react-native";
import Toast from "react-native-toast-message";
import { PersistGate } from "redux-persist/integration/react";
import { useAppDispatch, useAppSelector } from "@/store/custom";

// Prevent the splash screen from auto-hiding before asset loading is complete.
SplashScreen.preventAutoHideAsync();

export default function RootLayout() {
  return (
    <ReduxProvider store={store}>
      <PersistGate persistor={persistor}>
        <AuthGate>
          <PaperProvider theme={MD3LightTheme}>
            <App />
            <Toast />
          </PaperProvider>
        </AuthGate>
      </PersistGate>
    </ReduxProvider>
  );
}

function AuthGate({ children }: { children: React.ReactNode }) {
  const dispatch = useAppDispatch();
  const isResolved = useAppSelector(
    (state) => state.auth.type !== "unresolved",
  );

  useEffect(() => {
    const loadAuthState = async () => {
      try {
        const token = await SecureStore.getItemAsync(SESSION_TOKEN_KEY);
        if (token) {
          dispatch(validateStoredToken(token));
        } else {
          dispatch(setUnauthenticated());
        }
      } catch (error) {
        console.error("Error loading auth state:", error);
      }
    };

    loadAuthState();
  }, [dispatch]);

  useEffect(() => {
    if (isResolved) {
      SplashScreen.hideAsync();
    }
  }, [isResolved]);

  if (isResolved) {
    return children;
  }
  {
    return null;
  }
}

function App() {
  const theme = useTheme();

  return (
    <Stack
      screenOptions={{
        header: (props) => <Header {...props} />,
        contentStyle: {
          backgroundColor: theme.colors.background,
        },
      }}
      layout={({ children }) => (
        <View style={{ flex: 1 }}>
          <StatusBar backgroundColor={theme.colors.surface} />
          {children}
        </View>
      )}
    >
      <Stack.Screen name="index" options={{ headerShown: false }} />
      <Stack.Screen name="login" options={{ headerShown: false }} />
      <Stack.Screen name="(tabs)" options={{ headerShown: false }} />
    </Stack>
  );
}
