import {create} from 'zustand';

interface Client {
  id: string;
  name?: string;
  clientSecret?: string;
  createdAt?: Date;
}

interface AccountState {
  client: Client | null;
  clients: Client[];
  setClient: (client: Client) => void;
  setClients: (clients: Client[]) => void;
}

export const useClientStore = create<AccountState>((set) => ({
  client: null,
  clients: [],
  setClient: (client: Client) => set({client}),
  setClients: (clients: Client[]) => set({clients}),
}));