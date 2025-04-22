import {useClientStore} from "@/store/clients.ts";
import {useEffect} from "react";
import {usePoolStore} from "@/store/pool.ts";
import {authbase} from "@/api/client.ts";
import dayjs from "dayjs";

export const useListClients = () => {
  const activePool = usePoolStore(state => state.activePool);

  useEffect(() => {
    if (!activePool) {
      return;
    }

    authbase.client.listClients({
      pool_id: activePool?.id,
    }).then((res) => {
      const {data} = res;
      const clients = data.clients?.map((client) => ({
        id: client.id!,
        name: client.name!,
        createdAt: dayjs(client.createdAt).toDate() ?? undefined,
        account: {
          id: client.CreatedByUser?.id ?? '',
          name: client.CreatedByUser?.name,
          email: client.CreatedByUser?.email ?? '',
        }
      })) || [];
      useClientStore.getState().setClients(clients);
    })
  }, [activePool]);
}

export const useGetAccount = () => {
}