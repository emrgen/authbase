import {create} from "zustand";

interface ActiveSidebarItem {
  title: string;
  url: string;
  icon: React.ElementType;
}

interface AppStore {
  isLoading: boolean;
  setLoading: (loading: boolean) => void;
  activeSidebarItem?: ActiveSidebarItem | null;
  setActiveSidebarItem: (item: ActiveSidebarItem | null) => void;
}

export const useAppStore = create<AppStore>((set) => ({
  isLoading: false,
  setLoading: (loading) => set({isLoading: loading}),
  activeSidebarItem: null,
  setActiveSidebarItem: (item) => set({activeSidebarItem: item}),
}));