import { Ionicons } from "@expo/vector-icons";
import { router, Tabs } from "expo-router";
import { Platform, Pressable } from "react-native";
import { useAuth } from "@/context/auth";
export const unstable_settings = {
  initialRouteName: "index",
};

export default function TabLayout() {
  const { user } = useAuth();

  return (
    <Tabs>
      <Tabs.Screen
        name="(outings)/index"
        options={{
          title: "Outings",
          ...(!user && {
            href: null,
          }),
        }}
      />
      <Tabs.Screen
        name="(outings)/outings/[id]"
        options={{
          href: null,
        }}
      />
    </Tabs>
  );
}
