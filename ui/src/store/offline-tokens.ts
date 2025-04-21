import {create} from "zustand";

interface OfflineToken {
  id: string;
  name: string;
  token: string;
  createdBy: string;
  createdAt: number; // Unix timestamp in seconds
  expiresAt: number; // Unix timestamp in seconds
}

export interface OfflineTokenStore {
  tokens: OfflineToken[];
  addToken: (token: OfflineToken) => void;
  removeToken: (token: string) => void;
  clearTokens: () => void;
}

export const useOfflineTokenStore = create<OfflineTokenStore>((set) => ({
  tokens: [],
  addToken: (token) =>
    set((state) => ({
      tokens: [...state.tokens, token],
    })),
  removeToken: (id: string) =>
    set((state) => ({
      tokens: state.tokens.filter((t) => t.id !== id),
    })),
  clearTokens: () => set({tokens: []}),
}));