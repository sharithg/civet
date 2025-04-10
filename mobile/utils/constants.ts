import { Platform } from "react-native";

console.log({ os: Platform.OS });

export const API_URL =
  Platform.OS === "ios"
    ? process.env.EXPO_PUBLIC_API_URL
    : process.env.EXPO_PUBLIC_LOCAL_API_URL;
