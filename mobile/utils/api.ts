import axios, { AxiosRequestConfig } from "axios";
import { Platform } from "react-native";
import { tokenCache } from "./cache";
import { API_URL } from "./constants";

export const authFetch = async <T = any>(
  path: string,
  config: AxiosRequestConfig = { method: "GET" }
): Promise<T> => {
  const token = await tokenCache.getToken("accessToken");

  const defaultHeaders = {
    ...(token ? { Authorization: `Bearer ${token}` } : {}),
    Platform: Platform.OS,
  };

  const finalConfig: AxiosRequestConfig = {
    ...config,
    headers: {
      ...defaultHeaders,
      ...(config.headers || {}),
    },
    withCredentials: true,
  };

  const normalizedPath = path.startsWith("/") ? path.slice(1) : path;
  const url = `${API_URL!.replace(/\/$/, "")}/${normalizedPath}`;

  const res = await axios({
    url,
    ...finalConfig,
  });
  return res.data;
};
