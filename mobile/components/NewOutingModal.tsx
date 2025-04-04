import { useMutation } from "@tanstack/react-query";
import {
  Modal,
  View,
  Text,
  TextInput,
  TouchableOpacity,
  StyleSheet,
} from "react-native";
import { router } from "expo-router";
import { authFetch } from "@/utils/api";

const newOutingRequest = async (name: string) => {
  const data = await authFetch<{ id: string }>("outing", {
    method: "POST",
    data: { name },
  });
  return data;
};

export default function NewOutingModal({
  modalVisible,
  setModalVisible,
  newOutingName,
  setNewOutingName,
  refetch,
}: {
  modalVisible: boolean;
  setModalVisible: (v: boolean) => void;
  setNewOutingName: (t: string) => void;
  newOutingName: string;
  refetch: () => Promise<void>;
}) {
  const { mutateAsync } = useMutation<
    { id: string },
    unknown,
    { name: string }
  >({
    mutationFn: (input) => newOutingRequest(input.name),
  });

  return (
    <Modal
      visible={modalVisible}
      animationType="slide"
      transparent
      onRequestClose={() => setModalVisible(false)}
    >
      <View style={styles.modalOverlay}>
        <View style={styles.modalContainer}>
          <Text style={styles.modalTitle}>New Outing</Text>
          <TextInput
            style={styles.input}
            placeholder="Enter outing name"
            value={newOutingName}
            onChangeText={setNewOutingName}
          />
          <TouchableOpacity
            onPress={async () => {
              setModalVisible(false);
              const resp = await mutateAsync({ name: newOutingName });
              setNewOutingName("");
              router.navigate({
                pathname: "/dashboard/(outings)/outings/[id]",
                params: {
                  id: resp.id,
                  title: newOutingName,
                },
              });
              await refetch();
            }}
            style={styles.createButton}
          >
            <Text style={styles.createButtonText}>Create</Text>
          </TouchableOpacity>
          <TouchableOpacity onPress={() => setModalVisible(false)}>
            <Text style={{ textAlign: "center", marginTop: 10, color: "#888" }}>
              Cancel
            </Text>
          </TouchableOpacity>
        </View>
      </View>
    </Modal>
  );
}

const styles = StyleSheet.create({
  modalOverlay: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
    backgroundColor: "rgba(0,0,0,0.4)",
  },
  modalContainer: {
    backgroundColor: "white",
    width: "80%",
    padding: 20,
    borderRadius: 12,
    elevation: 5,
  },
  modalTitle: {
    fontSize: 18,
    fontWeight: "600",
    marginBottom: 12,
  },
  input: {
    borderWidth: 1,
    borderColor: "#ccc",
    borderRadius: 8,
    padding: 10,
    marginBottom: 16,
  },
  createButton: {
    backgroundColor: "#111",
    paddingVertical: 12,
    borderRadius: 8,
  },
  createButtonText: {
    color: "white",
    textAlign: "center",
    fontWeight: "600",
  },
});
