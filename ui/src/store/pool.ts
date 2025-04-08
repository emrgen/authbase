import {create} from 'zustand';

interface Pool {
  id: string;
  name: string;
  description?: string;
  createdAt?: Date;
}

interface PoolStore {
  pools: Pool[];
  activePool: Pool | null;
  setPools: (pools: Pool[]) => void;
  setActivePool: (pool: Pool | null) => void;
}

export const usePoolStore = create<PoolStore>((set) => ({
  pools: [],
  activePool: null,
  setPools: (pools: Pool[]) => {
    set({pools});
  },
  setActivePool: (pool: Pool | null) => {
    set({activePool: pool});
  },
}));