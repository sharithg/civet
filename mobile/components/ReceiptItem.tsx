import React from "react";
import { View, Text, ScrollView, StyleSheet } from "react-native";

type Item = {
  id: string;
  name: string;
  price: number;
  quantity: number;
};

type Receipt = {
  restaurant: string;
  address: string;
  server: string;
  order_number: string;
  order_type: string;
  sales_tax: number;
  total: number;
  items: Item[];
};

const ReceiptView: React.FC<{ receipt: Receipt }> = ({ receipt }) => {
  console.log({ receipt });
  return (
    <ScrollView contentContainerStyle={styles.container}>
      {/* Header Info */}
      <Text style={styles.address}>{receipt.address}</Text>

      <View style={styles.metaRow}>
        <Text>Server: {receipt.server}</Text>
        <Text>Order #: {receipt.order_number}</Text>
      </View>

      <Text style={styles.orderType}>Type: {receipt.order_type}</Text>

      {/* Items */}
      <View style={styles.sectionHeader}>
        <Text style={styles.sectionHeaderText}>Items</Text>
      </View>

      {receipt.items.map((item) => (
        <View key={item.id} style={styles.itemRow}>
          <View style={styles.itemLeft}>
            <Text style={styles.itemName}>{item.name}</Text>
            <Text style={styles.quantity}>x{item.quantity}</Text>
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
  },
  itemName: {
    fontSize: 14,
  },
  quantity: {
    fontSize: 12,
    color: "#888",
  },
  itemPrice: {
    fontSize: 14,
    fontWeight: "500",
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
});

export default ReceiptView;
