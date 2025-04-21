import {usePoolStore} from "@/store/pool.ts";
import {useEffect} from "react";

export const useListProviders = () => {
  const activePool = usePoolStore(state => state.activePool);

  useEffect(() => {
    if (!activePool) {
      return;
    }

  //   authbase.provider.listClients({
  //     pool_id: activePool?.id,
  //   }).then((res) => {
  //     const {data} = res;
  //     const clients = data.clients?.map((client) => ({
  //       id: client.id!,
  //       name: client.name!,
  //       clientSecret: client.clientSecret!,
  //       createdAt: dayjs(client.createdAt).toDate() ?? undefined,
  //     })) || [];
  //     useClientStore.getState().setClients(clients);
  //   })
  }, [activePool]);
}