import {
  DarkTheme,
  DefaultTheme,
  ThemeProvider,
} from "@react-navigation/native";
import { useFonts } from "expo-font";
import { Slot, Stack, Tabs } from "expo-router";
import * as SplashScreen from "expo-splash-screen";
import { StatusBar } from "expo-status-bar";
import { useEffect } from "react";
import "react-native-reanimated";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import Toast from "react-native-toast-message";
import { AuthProvider, useAuth } from "../context/auth";

const queryClient = new QueryClient();

// Prevent the splash screen from auto-hiding before asset loading is complete.
SplashScreen.preventAutoHideAsync();

export default function RootLayout() {
  // const colorScheme = useColorScheme();
  const [loaded] = useFonts({
    SpaceMono: require("../assets/fonts/SpaceMono-Regular.ttf"),
  });

  useEffect(() => {
    if (loaded) {
      SplashScreen.hideAsync();
    }
  }, [loaded]);

  if (!loaded) {
    return null;
  }

  return (
    <AuthProvider>
      <QueryClientProvider client={queryClient}>
        <ThemeProvider value={DefaultTheme}>
          <Tabs>
            <Tabs.Screen
              name="dashboard/(outings)"
              options={{
                title: "Home",
                headerShown: false,
              }}
            />
            <Tabs.Screen
              name="(auth)"
              options={{
                href: null,
              }}
            />
            <Tabs.Screen
              name="+not-found"
              options={{
                href: null,
              }}
            />
          </Tabs>
        </ThemeProvider>
        <Toast />
      </QueryClientProvider>
    </AuthProvider>
  );
}
