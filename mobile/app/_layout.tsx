import { Stack } from "expo-router";
import * as SplashScreen from "expo-splash-screen";
import { PaperProvider } from "react-native-paper";
import { useEffect } from "react";
import Header from "@/components/header";
import { Provider as ReduxProvider } from "react-redux";
import { store, useAppDispatch, useAppSelector } from "@/store";
import * as SecureStore from "expo-secure-store";
// Initialize localization.
import "@/components/i18n";
import { fetchMe, setUnauthenticated } from "@/components/auth/auth-slice";
import { SESSION_TOKEN_KEY } from "@/components/auth/constants";

// Prevent the splash screen from auto-hiding before asset loading is complete.
SplashScreen.preventAutoHideAsync();

export default function RootLayout() {
  return (
    <ReduxProvider store={store}>
      <AuthGate>
        <PaperProvider>
          <App />
        </PaperProvider>
      </AuthGate>
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
          dispatch(fetchMe(token));
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
  return (
    <Stack
      screenOptions={{
        header: (props) => <Header {...props} />,
      }}
    >
      <Stack.Screen name="index" options={{ headerShown: false }} />
      <Stack.Screen name="login" options={{ headerShown: false }} />
      <Stack.Screen name="(tabs)" options={{ headerShown: false }} />
    </Stack>
  );
}
