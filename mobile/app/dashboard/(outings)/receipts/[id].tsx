import React from "react";
import { SafeAreaView, ScrollView, StyleSheet } from "react-native";
import ReceiptView from "@/components/ReceiptItem";
import BackButton from "../../../../components/BackButton";

const receiptData = {
  id: "9f9d6c11-7afa-4601-91ed-87ab9ff9701d",
  total: 7.61,
  restaurant: "Taco Bell",
  address: "7230 Lawrence Pendleton Pike, IN 46226",
  opened: "0001-01-01 00:00:00.000000",
  order_number: "378752",
  order_type: "drive-thru",
  server: "DAJA G",
  sales_tax: 0.63,
  items: [
    {
      id: "060ef667-f6c5-4cd3-8004-34d3c83b1e5d",
      name: "Power Veg Bowl No Sour Cream",
      price: 4.99,
      quantity: 1,
    },
    {
      id: "4b200c31-e7dd-444f-90cd-f30a3c62ade5",
      name: "Rg Orange Crsh Fz",
      price: 1.99,
      quantity: 1,
    },
  ],
  fees: [],
};

export default function Receipts() {
  return (
    <>
      <BackButton title={receiptData.restaurant} />
      <SafeAreaView style={styles.container}>
        <ScrollView style={{ flex: 1, backgroundColor: "#f8f8f8" }}>
          <ReceiptView receipt={receiptData} />
        </ScrollView>
      </SafeAreaView>
    </>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    fontSize: 30,
  },
});
