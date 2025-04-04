import { Ionicons } from "@expo/vector-icons";
import { router, Stack } from "expo-router";
import { Platform } from "react-native";

export default function BackButton({ title }: { title: string }) {
  const backIcon = Platform.OS === "ios" ? "chevron-back" : "arrow-back-sharp";

  return (
    <Stack.Screen
      options={{
        headerShown: true,
        headerLeft: () => (
          <Ionicons
            name={backIcon}
            size={25}
            color="blue"
            onPress={() => router.back()}
          />
        ),
        title: title as string,
      }}
    />
  );
}
