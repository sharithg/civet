import { Receipt, Split } from "@/app/dashboard/(outings)/receipts/[id]";
import { authFetch } from "@/utils/api";
import { useDebouncedMutation } from "@/utils/hooks/useDebouncedMutation";
import { Friend, selectedOutingAtom } from "@/utils/state";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useAtomValue } from "jotai";
import React, { useEffect, useState } from "react";
import {
  View,
  Text,
  ScrollView,
  StyleSheet,
  TextInput,
  TouchableOpacity,
  Image,
} from "react-native";

type NewFriend = { receiptId: string; name: string };
type NewFriendResponse = { friendId: string };

type CreateSplitItem = {
  friend_id: string;
  item_id: string;
  quantity: number;
};

type CreateSplitInput = {
  receipt_id: string;
  items: CreateSplitItem[];
};

const fetchFriends = async (receiptId: string) => {
  const result = authFetch<Friend[]>(`receipt/${receiptId}/friends`);
  return result || [];
};

const addNewFriend = async (req: NewFriend) => {
  console.log({ req });
  const result = await authFetch<NewFriendResponse>(`receipt/friends`, {
    method: "POST",
    data: JSON.stringify({
      name: req.name,
      receipt_id: req.receiptId,
      user_id: null,
    }),
  });
  return result || [];
};

const createSplit = async (req: CreateSplitInput) => {
  const result = await authFetch<null>(`receipt/friends/split`, {
    method: "POST",
    data: JSON.stringify(req),
  });
  return result;
};

type SplitState = Record<string, string[]>;

const splitToCreateSplit = (
  split: SplitState,
  receiptId: string
): CreateSplitInput => {
  const items: CreateSplitItem[] = [];

  Object.entries(split).forEach(([itemId, friendIds]) => {
    friendIds.forEach((friendId) => {
      items.push({
        item_id: itemId,
        friend_id: friendId,
        quantity: 1,
      });
    });
  });

  return {
    receipt_id: receiptId,
    items,
  };
};

const aggSplits = (splits: Split[]): Record<string, string[]> => {
  const splitMap: Record<string, string[]> = {};

  splits.forEach((split) => {
    const { order_item_id, friend_id } = split;

    if (!splitMap[order_item_id]) {
      splitMap[order_item_id] = [];
    }

    splitMap[order_item_id].push(friend_id);
  });

  return splitMap;
};

const ReceiptView: React.FC<{ receipt: Receipt }> = ({ receipt }) => {
  const [splitWith, setSplitWith] = useState<SplitState>(
    aggSplits(receipt.splits)
  );
  const [newFriend, setNewFriend] = useState("");
  const selectedOuting = useAtomValue(selectedOutingAtom);

  const { data, refetch } = useQuery({
    queryFn: () => fetchFriends(receipt.id),
    queryKey: ["friends", receipt.id],
  });
  const { mutateAsync } = useMutation<NewFriendResponse, unknown, NewFriend>({
    mutationFn: (data) => addNewFriend(data),
  });

  const { mutate: mutateSplit } = useDebouncedMutation<null, CreateSplitInput>({
    mutationFn: createSplit,
  });

  const [existingFriends, setExistingFriends] = useState<Friend[]>(data || []);

  useEffect(() => {
    if (data) {
      setExistingFriends(data);
    }
  }, [data]);

  const toggleFriend = async (itemId: string, friendId: string) => {
    setSplitWith((prev) => {
      const current = prev[itemId] || [];
      const updated = current.includes(friendId)
        ? current.filter((id) => id !== friendId)
        : [...current, friendId];

      const newSplitWith = {
        ...prev,
        [itemId]: updated,
      };

      const split = splitToCreateSplit(newSplitWith, receipt.id);
      mutateSplit(split);

      return newSplitWith;
    });
  };

  const handleAddFriend = async () => {
    const trimmed = newFriend.trim();
    if (!trimmed) return;
    await mutateAsync({
      name: trimmed,
      receiptId: receipt.id,
    });
    await refetch();
    setNewFriend("");
  };

  const getFriendNameById = (id: string) =>
    existingFriends.find((f) => f.id === id)?.name ?? id;

  console.log("splitWith", splitWith);
  return (
    <ScrollView contentContainerStyle={styles.container}>
      {/* Header Info */}
      <Text style={styles.address}>{receipt.address}</Text>
      <View style={styles.metaRow}>
        <Text>Server: {receipt.server}</Text>
        <Text>Order #: {receipt.order_number}</Text>
      </View>
      <Text style={styles.orderType}>Type: {receipt.order_type}</Text>

      <View style={styles.globalAddFriendContainer}>
        <TextInput
          placeholder="New friend"
          value={newFriend}
          onChangeText={setNewFriend}
          style={styles.friendInput}
        />
        <TouchableOpacity style={styles.addButton} onPress={handleAddFriend}>
          <Text style={styles.addButtonText}>Add</Text>
        </TouchableOpacity>
      </View>

      {/* Items */}
      <View style={styles.sectionHeader}>
        <Text style={styles.sectionHeaderText}>Items</Text>
      </View>

      {receipt.items.map((item) => (
        <View key={item.id} style={styles.itemRow}>
          <View style={styles.itemLeft}>
            <View style={{ flex: 1, flexDirection: "row" }}>
              <Text style={styles.quantity}>{`${item.quantity}  `}</Text>
              <Text style={styles.itemName}>{item.name}</Text>
            </View>

            <View style={styles.friendsList}>
              {existingFriends.map((friend) => (
                <TouchableOpacity
                  key={friend.id}
                  style={[
                    styles.friendButton,
                    splitWith[item.id]?.includes(friend.id) &&
                      styles.friendSelected,
                  ]}
                  onPress={() => toggleFriend(item.id, friend.id)}
                >
                  <Text style={styles.friendText}>{friend.name}</Text>
                </TouchableOpacity>
              ))}
            </View>

            {splitWith[item.id]?.length > 1 && (
              <Text style={styles.splitLabel}>
                Split between:{" "}
                {splitWith[item.id].map(getFriendNameById).join(", ")}
              </Text>
            )}
            {splitWith[item.id]?.length === 1 && (
              <Text style={styles.splitLabel}>
                Assigned to: {getFriendNameById(splitWith[item.id][0])}
              </Text>
            )}
          </View>

          <Text style={styles.itemPrice}>
            ${(item.price * item.quantity).toFixed(2)}
          </Text>
        </View>
      ))}

      {/* Totals */}
      <View style={styles.totals}>
        <View style={styles.totalRow}>
          <Text>Tax</Text>
          <Text>${receipt.sales_tax.toFixed(2)}</Text>
        </View>

        <Text style={styles.taxNote}>
          Tax split equally between all assigned friends.
        </Text>

        <View style={styles.totalRow}>
          <Text style={styles.totalLabel}>Total</Text>
          <Text style={styles.totalAmount}>${receipt.total.toFixed(2)}</Text>
        </View>
      </View>
    </ScrollView>
  );
};

const styles = StyleSheet.create({
  container: {
    padding: 16,
    backgroundColor: "#fff",
    borderRadius: 8,
  },
  globalAddFriendContainer: {
    flexDirection: "row",
    alignItems: "center",
    marginBottom: 16,
    gap: 8,
  },
  address: {
    fontSize: 14,
    color: "#666",
    marginBottom: 8,
  },
  metaRow: {
    flexDirection: "row",
    justifyContent: "space-between",
    marginBottom: 4,
  },
  orderType: {
    fontSize: 14,
    marginBottom: 12,
  },
  sectionHeader: {
    borderBottomWidth: 1,
    borderColor: "#eee",
    marginBottom: 8,
  },
  sectionHeaderText: {
    fontSize: 16,
    fontWeight: "bold",
  },
  itemRow: {
    flexDirection: "row",
    justifyContent: "space-between",
    paddingVertical: 6,
    borderBottomWidth: 1,
    borderColor: "#f0f0f0",
  },
  itemLeft: {
    flex: 1,
    paddingRight: 8,
  },
  itemName: {
    fontSize: 14,
  },
  quantity: {
    fontSize: 14,
    color: "#888",
  },
  itemPrice: {
    fontSize: 14,
    fontWeight: "500",
  },
  friendsList: {
    flexDirection: "row",
    flexWrap: "wrap",
    marginTop: 4,
    gap: 4,
  },
  friendButton: {
    backgroundColor: "#eee",
    borderRadius: 8,
    paddingVertical: 4,
    paddingHorizontal: 8,
    marginRight: 6,
    marginTop: 4,
  },
  friendSelected: {
    backgroundColor: "#cce5ff",
  },
  friendText: {
    fontSize: 12,
  },
  friendInput: {
    borderWidth: 1,
    borderColor: "#ddd",
    borderRadius: 6,
    padding: 6,
    fontSize: 13,
    marginTop: 4,
    width: "80%",
  },
  totals: {
    marginTop: 12,
  },
  totalRow: {
    flexDirection: "row",
    justifyContent: "space-between",
    paddingVertical: 4,
  },
  totalLabel: {
    fontWeight: "bold",
    fontSize: 16,
  },
  totalAmount: {
    fontWeight: "bold",
    fontSize: 16,
  },
  newFriendContainer: {
    flexDirection: "row",
    alignItems: "center",
    marginTop: 6,
    gap: 8,
  },
  addButton: {
    backgroundColor: "#007bff",
    paddingVertical: 6,
    paddingHorizontal: 12,
    borderRadius: 6,
    zIndex: 2,
  },
  addButtonText: {
    color: "#fff",
    fontWeight: "500",
    fontSize: 13,
  },
  taxNote: {
    fontSize: 12,
    color: "#888",
    fontStyle: "italic",
    marginTop: -4,
    marginBottom: 8,
  },
  splitLabel: {
    fontSize: 12,
    marginTop: 4,
    color: "#555",
  },
});

export default ReceiptView;
