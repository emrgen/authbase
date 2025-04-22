import {create} from 'zustand';
import { Account } from './account';

interface Client {
  id: string;
  name?: string;
  clientSecret?: string;
  createdAt?: Date;
  account?: Account;
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