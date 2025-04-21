interface Provider {
  name: string;
  clientId: string;
  clientSecret: string;
}


import {create} from 'zustand';

interface ProviderStore {
  providers: Provider[];
  setProviders: (providers: Provider[]) => void;
}

export const useProviderStore = create<ProviderStore>((set) => ({
  providers: [],
  setProviders: (providers) => set({providers}),
}));