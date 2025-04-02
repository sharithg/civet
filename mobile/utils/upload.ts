import * as ImagePicker from "expo-image-picker";
import { API_URL } from "./constants";

export const pickDocument = async () => {
  let result = await ImagePicker.launchImageLibraryAsync({
    mediaTypes: ["images"],
    allowsEditing: true,
    aspect: [4, 3],
    quality: 1,
  });

  console.log(result);

  if (!result.canceled) {
    return { uri: result.assets[0].uri, fileName: result.assets[0].fileName };
  }

  return null;
};

export const uploadImage = async (
  uri: string,
  fileName: string,
  outingId: string
) => {
  const match = /\.(\w+)$/.exec(fileName);

  const response = await fetch(uri);
  const blob = await response.blob();

  const formData = new FormData();
  formData.append("file", blob, fileName);

  try {
    const uploadResponse = await fetch(`${API_URL}/api/v1/receipt/upload`, {
      method: "POST",
      body: formData,
      headers: {
        outingId,
      },
    });

    const data = await uploadResponse.json();
    console.log(data);
  } catch (error) {
    console.error("Upload failed:", error);
  }
};
