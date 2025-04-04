import * as ImagePicker from "expo-image-picker";
import { API_URL } from "./constants";
import { tokenCache } from "./cache";

export const pickDocument = async () => {
  let result = await ImagePicker.launchImageLibraryAsync({
    mediaTypes: ["images"],
    allowsEditing: true,
    aspect: [4, 3],
    quality: 1,
  });

  console.log(result);

  if (!result.canceled) {
    return result;
  }

  return null;
};

function formDataFromImagePicker(result: ImagePicker.ImagePickerSuccessResult) {
  const formData = new FormData();

  for (const index in result.assets) {
    const asset = result.assets[index];

    console.log(`photo.${index}`);

    // @ts-expect-error: special react native format for form data
    formData.append(`photo.${index}`, {
      uri: asset.uri,
      name: asset.fileName ?? asset.uri.split("/").pop(),
      type: asset.mimeType,
    });

    if (asset.exif) {
      formData.append(`exif.${index}`, JSON.stringify(asset.exif));
    }
  }

  return formData;
}

export const uploadImage = async (
  result: ImagePicker.ImagePickerSuccessResult,
  outingId: string
) => {
  try {
    const token = await tokenCache.getToken("accessToken");
    const uploadResponse = await fetch(`${API_URL}/receipt/upload`, {
      method: "POST",
      body: formDataFromImagePicker(result),
      headers: {
        Accept: "application/json",
        outingId,
        Authorization: `Bearer ${token}`,
      },
      credentials: "include",
    });

    const data = await uploadResponse.json();
    console.log(data);
  } catch (error) {
    console.error("Upload failed:", error);
  }
};
