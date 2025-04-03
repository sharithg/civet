import { Ionicons } from "@expo/vector-icons";
import { router, Stack, Tabs } from "expo-router";
import { Platform, Pressable } from "react-native";
import { useAuth } from "@/context/auth";
import { useEffect } from "react";
export const unstable_settings = {
  initialRouteName: "index",
};

export default function Layout() {
  const { user, isLoading } = useAuth();

  useEffect(() => {
    if (user && !isLoading) {
      router.push("/dashboard");
    }
  }, [user, isLoading]);

  return (
    <Stack>
      <Stack.Screen
        name="index"
        options={{
          title: "Login",
        }}
      />
    </Stack>
  );
}
