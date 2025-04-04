import React from "react";
import { SafeAreaView, ScrollView, StyleSheet } from "react-native";
import ReceiptView from "@/components/ReceiptItem";
import BackButton from "../../../../components/BackButton";
import { authFetch } from "@/utils/api";
import { useLocalSearchParams } from "expo-router";
import { useQuery } from "@tanstack/react-query";

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

export interface ReceiptItem {
  id: string;
  receipt_id: string;
  name: string;
  price: number;
  quantity: number;
}

export interface Receipt {
  id: string;
  total: number;
  restaurant: string;
  address: string;
  opened: string; // ISO timestamp string
  order_number: string;
  order_type: string;
  payment_tip: number | null;
  payment_amount_paid: number | null;
  table_number: string;
  copy: string;
  server: string;
  sales_tax: number;
  items: ReceiptItem[];
  fees: any[];
}

const fetchReceipt = async (id: string) => {
  const result = authFetch<Receipt>(`receipt/item/${id}`);
  return result || [];
};

export default function Receipts() {
  const { id } = useLocalSearchParams();

  const receiptId = id as string;

  const { data } = useQuery({
    queryFn: () => fetchReceipt(receiptId),
    queryKey: [receiptId],
  });

  return (
    <>
      <BackButton title={receiptData.restaurant} />
      <SafeAreaView style={styles.container}>
        <ScrollView style={{ flex: 1, backgroundColor: "#f8f8f8" }}>
          {data && <ReceiptView receipt={data} />}
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
