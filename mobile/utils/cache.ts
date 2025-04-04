import * as SecureStore from "expo-secure-store";
import { Platform } from "react-native";
import Cookies from "js-cookie";

const createTokenCache = () => {
  return {
    getToken: async (key: string) => {
      if (Platform.OS === "web") {
        return Cookies.get(key) ?? null;
      }
      try {
        return await SecureStore.getItemAsync(key);
      } catch {
        return null;
      }
    },
    saveToken: async (key: string, token: string) => {
      if (Platform.OS === "web") {
        Cookies.set(key, token, { secure: true, sameSite: "Lax" });
      } else {
        await SecureStore.setItemAsync(key, token);
      }
    },
    deleteToken: async (key: string) => {
      if (Platform.OS === "web") {
        Cookies.remove(key);
      } else {
        await SecureStore.deleteItemAsync(key);
      }
    },
  };
};

// SecureStore is not supported on the web and we use cookies instead
export const tokenCache = createTokenCache();
