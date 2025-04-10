import * as ImagePicker from "expo-image-picker";
import { API_URL } from "./constants";
import { tokenCache } from "./cache";
import { Platform } from "react-native";

function dataURItoBlob(dataURI: string) {
  // convert base64/URLEncoded data component to raw binary data held in a string
  var byteString;
  if (dataURI.split(",")[0].indexOf("base64") >= 0)
    byteString = atob(dataURI.split(",")[1]);
  else byteString = unescape(dataURI.split(",")[1]);

  // separate out the mime component
  var mimeString = dataURI.split(",")[0].split(":")[1].split(";")[0];

  // write the bytes of the string to a typed array
  var ia = new Uint8Array(byteString.length);
  for (var i = 0; i < byteString.length; i++) {
    ia[i] = byteString.charCodeAt(i);
  }

  return new Blob([ia], { type: mimeString });
}

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

function formDataFromImagePickerWeb(
  result: ImagePicker.ImagePickerSuccessResult
) {
  const formData = new FormData();

  for (const index in result.assets) {
    const asset = result.assets[index];

    console.log(`photo.${index}`);

    var blob = dataURItoBlob(asset.uri);

    formData.append(`photo.${index}`, blob, asset.fileName ?? asset.uri);

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
      body:
        Platform.OS === "web"
          ? formDataFromImagePickerWeb(result)
          : formDataFromImagePicker(result),
      headers: {
        Accept: "application/json",
        outingId,
        Authorization: `Bearer ${token}`,
        platform: Platform.OS,
      },
      credentials: "include",
    });

    const data = await uploadResponse.json();
    console.log(data);
  } catch (error) {
    console.error("Upload failed:", error);
  }
};
