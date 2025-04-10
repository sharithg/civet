import React from "react";
import { View, Text, FlatList, StyleSheet } from "react-native";
import { User } from "lucide-react-native";
import { authFetch } from "@/utils/api";
import { useQuery } from "@tanstack/react-query";
import LoadingProgress from "./LoadingProgress";

interface FriendData {
  name: string;
  subtotal: number;
  tax_portion: number;
  total_owed: number;
}

interface FriendsListProps {
  outingId: string;
}

const fetchFriends = async (outingId: string) => {
  const result = authFetch<FriendData[]>(`outing/${outingId}/friends`);
  return result || [];
};

export default function FriendsList({ outingId }: FriendsListProps) {
  const { data, isLoading } = useQuery({
    queryFn: () => fetchFriends(outingId),
    queryKey: ["friends", outingId],
  });

  if (isLoading) {
    return <LoadingProgress />;
  }

  if (data?.length === 0) {
    return (
      <View style={styles.emptyContainer}>
        <User size={48} color="#D1D5DB" />
        <Text style={styles.emptyTitle}>No friends added yet</Text>
        <Text style={styles.emptySubtitle}>
          Add friends to split the bill with
        </Text>
      </View>
    );
  }

  return (
    <FlatList
      data={data || []}
      keyExtractor={(_, index) => index.toString()}
      contentContainerStyle={{ padding: 16 }}
      ItemSeparatorComponent={() => <View style={{ height: 12 }} />}
      renderItem={({ item }) => (
        <View style={styles.card}>
          <View style={styles.header}>
            <View style={styles.avatarRow}>
              <View style={styles.avatar}>
                <Text style={styles.avatarText}>
                  {item.name.charAt(0).toUpperCase()}
                </Text>
              </View>
              <Text style={styles.name}>{item.name}</Text>
            </View>
            <Text style={styles.totalOwed}>${item.total_owed.toFixed(2)}</Text>
          </View>
        </View>
      )}
    />
  );
}

const styles = StyleSheet.create({
  emptyContainer: {
    alignItems: "center",
    paddingVertical: 40,
  },
  emptyTitle: {
    marginTop: 8,
    color: "#6B7280",
    fontSize: 16,
  },
  emptySubtitle: {
    color: "#9CA3AF",
    fontSize: 13,
  },
  card: {
    backgroundColor: "#FFFFFF",
    padding: 16,
    borderRadius: 12,
    elevation: 2,
    shadowColor: "#000",
    shadowOpacity: 0.05,
    shadowOffset: { width: 0, height: 2 },
    shadowRadius: 4,
  },
  header: {
    flexDirection: "row",
    justifyContent: "space-between",
    marginBottom: 12,
    alignItems: "center",
  },
  avatarRow: {
    flexDirection: "row",
    alignItems: "center",
  },
  avatar: {
    width: 40,
    height: 40,
    borderRadius: 20,
    backgroundColor: "#DBEAFE",
    alignItems: "center",
    justifyContent: "center",
    marginRight: 12,
  },
  avatarText: {
    color: "#2563EB",
    fontWeight: "600",
    fontSize: 16,
  },
  name: {
    fontSize: 16,
    fontWeight: "500",
  },
  totalOwed: {
    fontSize: 18,
    fontWeight: "700",
  },
  details: {
    gap: 4,
  },
  row: {
    flexDirection: "row",
    justifyContent: "space-between",
  },
  label: {
    color: "#6B7280",
    fontSize: 14,
  },
  value: {
    fontSize: 14,
  },
  bold: {
    fontWeight: "600",
  },
});
