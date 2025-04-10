import React from "react";
import { TouchableOpacity, StyleSheet, View, Text } from "react-native";
import { Plus } from "lucide-react-native";

interface FloatingActionButtonProps {
  onPress: () => void;
}

export default function FloatingActionButton({
  onPress,
}: FloatingActionButtonProps) {
  return (
    <TouchableOpacity style={styles.fab} onPress={onPress} activeOpacity={0.8}>
      <View style={styles.iconContainer}>
        <Plus size={24} color="#fff" />
      </View>
    </TouchableOpacity>
  );
}

const styles = StyleSheet.create({
  fab: {
    position: "absolute",
    bottom: 32,
    right: 24,
    zIndex: 10,
    elevation: 5,
    backgroundColor: "#2563EB", // Tailwind "blue-600"
    width: 56,
    height: 56,
    borderRadius: 28,
    justifyContent: "center",
    alignItems: "center",
    shadowColor: "#000",
    shadowOpacity: 0.2,
    shadowOffset: { width: 0, height: 3 },
    shadowRadius: 4,
  },
  iconContainer: {
    alignItems: "center",
    justifyContent: "center",
  },
});
