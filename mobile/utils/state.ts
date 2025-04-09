import { atom } from "jotai";

export type Friend = {
  id: string;
  name: string;
};

export type OutingData = {
  id: string;
  name: string;
  total_receipts: number;
  created_at: string;
  friends: Friend[];
  status: string;
};

export const selectedOutingAtom = atom<OutingData | undefined>();
