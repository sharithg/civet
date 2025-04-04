import React, { useState } from "react";
import {
  View,
  Text,
  StyleSheet,
  TouchableOpacity,
  FlatList,
  Platform,
} from "react-native";
import { Ionicons, MaterialCommunityIcons } from "@expo/vector-icons";
import { router, Stack, useLocalSearchParams } from "expo-router";
import { useQuery } from "@tanstack/react-query";
import { pickDocument, uploadImage } from "@/utils/upload";
import Toast from "react-native-toast-message";
import { authFetch } from "@/utils/api";
import BackButton from "../../../../components/BackButton";

type ReceiptData = {
  id: string;
  restaurant: string;
  order_count: number;
  total: number;
};

const fetchReceipts = async (id: string) => {
  const result = await authFetch<ReceiptData[]>(`outing/${id}/receipts`);
  return result;
};

export default function OutingDetailScreen() {
  const [selectedTab, setSelectedTab] = useState("Receipts");
  const backIcon = Platform.OS === "ios" ? "chevron-back" : "arrow-back-sharp";
  const { id, title } = useLocalSearchParams();

  const receiptId = id as string;

  const { data, refetch } = useQuery({
    queryFn: () => fetchReceipts(receiptId),
    queryKey: [receiptId],
  });

  return (
    <>
      <BackButton title={title as string} />
      <View style={styles.container}>
        {/* Tab Control */}
        <View style={styles.tabs}>
          {["Receipts", "Friends", "Split Bill"].map((tab) => (
            <TouchableOpacity
              key={tab}
              onPress={() => setSelectedTab(tab)}
              style={[
                styles.tabButton,
                selectedTab === tab && styles.activeTabButton,
              ]}
            >
              <Text
                style={[
                  styles.tabText,
                  selectedTab === tab && styles.activeTabText,
                ]}
              >
                {tab}
              </Text>
            </TouchableOpacity>
          ))}
        </View>

        {/* Action Button */}
        <TouchableOpacity
          style={styles.scanButton}
          onPress={async () => {
            const errorToast = () =>
              Toast.show({
                type: "error",
                text1: "error picking image",
                position: "bottom",
              });
            try {
              const result = await pickDocument();

              if (!result) {
                errorToast();
                return;
              }

              await uploadImage(result, id as string);
              await refetch();
            } catch (e) {
              errorToast();
            }
          }}
        >
          <Ionicons name="camera-outline" size={18} color="white" />
          <Text style={styles.scanText}>Scan New Receipt</Text>
        </TouchableOpacity>

        {/* Receipt List */}
        {selectedTab === "Receipts" && (
          <FlatList
            data={data ?? []}
            keyExtractor={(item) => item.id}
            contentContainerStyle={{ paddingBottom: 20 }}
            renderItem={({ item }) => (
              <TouchableOpacity
                style={styles.receiptCard}
                onPress={() =>
                  router.navigate({
                    pathname: `/dashboard/(outings)/receipts/[id]`,
                    params: {
                      id: item.id,
                    },
                  })
                }
              >
                <View style={styles.cardLeft}>
                  <MaterialCommunityIcons name="cash" size={28} color="#555" />
                  <View style={{ marginLeft: 12 }}>
                    <Text style={styles.receiptTitle}>{item.restaurant}</Text>
                    <Text style={styles.receiptSubtitle}>
                      {item.order_count} items
                    </Text>
                  </View>
                </View>
                <Text style={styles.amount}>${item.total.toFixed(2)}</Text>
              </TouchableOpacity>
            )}
          />
        )}
      </View>
    </>
  );
}

const styles = StyleSheet.create({
  container: { flex: 1, padding: 16, backgroundColor: "#f9f9f9" },
  header: {
    flexDirection: "row",
    alignItems: "center",
    gap: 12,
    marginBottom: 16,
  },
  title: { fontSize: 18, fontWeight: "600" },
  date: { color: "#888", fontSize: 14 },
  tabs: {
    flexDirection: "row",
    backgroundColor: "#f0f0f0",
    borderRadius: 8,
    marginBottom: 16,
  },
  tabButton: {
    flex: 1,
    paddingVertical: 10,
    alignItems: "center",
    borderRadius: 8,
  },
  activeTabButton: {
    backgroundColor: "white",
    borderWidth: 1,
    borderColor: "#aaa",
  },
  tabText: { color: "#888", fontSize: 14 },
  activeTabText: { color: "#111", fontWeight: "600" },
  scanButton: {
    flexDirection: "row",
    justifyContent: "center",
    alignItems: "center",
    backgroundColor: "#111",
    padding: 14,
    borderRadius: 8,
    marginBottom: 16,
  },
  scanText: {
    color: "white",
    marginLeft: 8,
    fontWeight: "600",
    fontSize: 14,
  },
  receiptCard: {
    backgroundColor: "white",
    flexDirection: "row",
    justifyContent: "space-between",
    alignItems: "center",
    padding: 16,
    marginBottom: 12,
    borderRadius: 12,
    shadowColor: "#000",
    shadowOpacity: 0.02,
    shadowOffset: { width: 0, height: 1 },
    shadowRadius: 2,
    elevation: 1,
  },
  cardLeft: { flexDirection: "row", alignItems: "center" },
  receiptTitle: { fontSize: 16, fontWeight: "500" },
  receiptSubtitle: { fontSize: 13, color: "#777" },
  amount: { fontWeight: "600", fontSize: 15 },
});
