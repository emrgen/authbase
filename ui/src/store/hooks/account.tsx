import {usePoolStore} from "../pool.ts";
import {authbase} from "../../api/client.ts";
import {useAccountStore} from "../account.ts";

export const useListAccounts = () => {
  const activePool = usePoolStore(state => state.activePool);

  useEffect(() => {
    if (!activePool) {
      return;
    }

    authbase.account.listAccounts({
      pool_id: activePool?.id,
    }).then((res) => {
      const {data} = res;
      const accounts = data.accounts?.map((account) => ({
        id: account.id!,
        username: account.username!,
        email: account.email!,
      })) || [];
      useAccountStore.getState().setAccounts(accounts);
    })
  }, [activePool]);
}

export const useGetAccount = () => {
}