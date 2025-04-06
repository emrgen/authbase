import {create} from 'zustand';

interface Account {
  id: string;
  name?: string;
  email: string;
  username?: string;
}

interface AccountState {
  account: Account | null;
  accounts: Account[];
  setAccount: (account: Account) => void;
  setAccounts: (accounts: Account[]) => void;
}

export const useAccountStore = create<AccountState>((set) => ({
  account: null,
  accounts: [],
  setAccount: (account: Account) => set({account}),
  setAccounts: (accounts: Account[]) => set({accounts}),
}));