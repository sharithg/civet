import { router } from "expo-router";
import {
  SafeAreaView,
  FlatList,
  StyleSheet,
  Text,
  View,
  ScrollView,
  TouchableOpacity,
  Platform,
} from "react-native";
import { SafeAreaProvider } from "react-native-safe-area-context";
import { Stack } from "expo-router";
import { Ionicons } from "@expo/vector-icons";
import NewOutingModal from "./NewOutingModal";
import { useState } from "react";
import axios from "axios";
import { useQuery } from "@tanstack/react-query";
import { format } from "date-fns";
import { API_URL } from "../utils/constants";
import * as AuthSession from "expo-auth-session";

type OutingData = {
  id: string;
  name: string;
  total_receipts: number;
  created_at: string;
  friends: string;
  status: string;
};

const fetchOutings = async () => {
  const result = await axios.get<OutingData[]>(`${API_URL}/outing`);
  return result.data || [];
};

export default function OutingsPage() {
  const [modalVisible, setModalVisible] = useState(false);
  const [newOutingName, setNewOutingName] = useState("");

  const { data } = useQuery({
    queryFn: fetchOutings,
    queryKey: ["outing"],
  });

  return (
    <SafeAreaProvider>
      <Stack.Screen
        options={{
          headerShown: true,
          headerRight: () => (
            <Ionicons
              name={"add-circle"}
              size={25}
              color="black"
              onPress={() => setModalVisible(true)}
              style={styles.addOutingButton}
            />
          ),
        }}
      />
      <NewOutingModal
        modalVisible={modalVisible}
        newOutingName={newOutingName}
        setModalVisible={(v) => setModalVisible(v)}
        setNewOutingName={(v) => setNewOutingName(v)}
      />
      <SafeAreaView style={styles.container}>
        <ScrollView style={{ padding: 10 }}>
          {(data ?? []).map((item) => (
            <TouchableOpacity
              key={item.id}
              onPress={() =>
                router.navigate({
                  pathname: "/dashboard/(outings)/outings/[id]",
                  params: {
                    id: item.id,
                    title: item.name,
                  },
                })
              }
            >
              <View
                style={{
                  backgroundColor: "#fff",
                  padding: 16,
                  marginBottom: 12,
                  borderRadius: 10,
                  shadowColor: "#000",
                  shadowOpacity: 0.05,
                  shadowRadius: 5,
                  elevation: 2,
                }}
              >
                <Text
                  style={{ fontSize: 18, fontWeight: "bold", marginBottom: 6 }}
                >
                  {item.name}
                </Text>

                <View
                  style={{
                    flexDirection: "row",
                    alignItems: "center",
                    marginBottom: 6,
                  }}
                >
                  <Text style={{ fontSize: 14, color: "#555" }}>
                    📅 {format(item.created_at, "PPPpp")}
                  </Text>
                </View>

                <View style={{ flexDirection: "row", marginBottom: 6 }}>
                  <Text
                    style={{ marginRight: 16, fontSize: 14, color: "#555" }}
                  >
                    👥 {item.friends} friends
                  </Text>
                  <Text style={{ fontSize: 14, color: "#555" }}>
                    💵 {item.total_receipts} receipts
                  </Text>
                </View>

                <View style={{ alignSelf: "flex-end" }}>
                  <Text
                    style={{
                      paddingHorizontal: 12,
                      paddingVertical: 4,
                      borderRadius: 20,
                      backgroundColor:
                        item.status === "Active" ? "#111" : "#eee",
                      color: item.status === "Active" ? "#fff" : "#111",
                      fontWeight: "bold",
                    }}
                  >
                    {item.status}
                  </Text>
                </View>
              </View>
            </TouchableOpacity>
          ))}
        </ScrollView>
      </SafeAreaView>
    </SafeAreaProvider>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    marginTop: 5,
    fontSize: 30,
  },
  item: {
    padding: 20,
    marginTop: 5,
    fontSize: 15,
  },
  addOutingButton: {
    marginRight: 10,
  },
});
