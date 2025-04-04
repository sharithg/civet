import { Ionicons } from "@expo/vector-icons";
import { router, Stack, Tabs } from "expo-router";
import { Platform, Pressable } from "react-native";
import { useAuth } from "@/context/auth";
export const unstable_settings = {
  initialRouteName: "index",
};

export default function TabLayout() {
  const { user } = useAuth();

  return <Stack />;
}
