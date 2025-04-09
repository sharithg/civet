import React from "react";
import {
  SafeAreaView,
  ScrollView,
  StyleSheet,
  View,
  Image,
} from "react-native";
import ReceiptView from "@/components/ReceiptItem";
import BackButton from "../../../../components/BackButton";
import { authFetch } from "@/utils/api";
import { useLocalSearchParams } from "expo-router";
import { useQuery } from "@tanstack/react-query";

export interface ReceiptItem {
  id: string;
  receipt_id: string;
  name: string;
  price: number;
  quantity: number;
}

export interface Split {
  id: string;
  friend_id: string;
  order_item_id: string;
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
  image_url: string;
  items: ReceiptItem[];
  fees: any[];
  splits: Split[];
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
      <BackButton title={data?.restaurant || ""} />
      <SafeAreaView style={styles.container}>
        <ScrollView style={{ flex: 1, backgroundColor: "#f8f8f8" }}>
          {data && (
            <>
              <ReceiptView receipt={data} />
              <Image
                style={styles.receptImage}
                source={{
                  uri: data.image_url,
                }}
              />
            </>
          )}
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
  receptImage: {
    width: "90%",
    height: undefined,
    aspectRatio: 2,
    alignSelf: "center",
    marginVertical: 20,
    resizeMode: "contain",
  },
});
