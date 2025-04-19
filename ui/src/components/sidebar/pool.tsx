import {authbase} from "@/api/client.ts";
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectLabel,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select"
import {usePoolStore} from "@/store/pool.ts";
import {useEffect} from "react";

export const SelectPool = () => {
  const pools = usePoolStore(state => state.pools);
  const activePool = usePoolStore(state => state.activePool);

  const setActivePool = usePoolStore(state => state.setActivePool);

  const handleChange = (poolId: string) => {
    const selectedPool = pools.find(pool => pool.id === poolId);
    if (selectedPool) {
      setActivePool(selectedPool);
    }
  };

  // if there are no pools, set the first pool as active
  useEffect(() => {
    if (pools.length > 0 && !pools.some(pool => pool.id === usePoolStore.getState().activePool?.id)) {
      const firstPool = pools[0];
      setActivePool(firstPool);
    }
  }, [pools, setActivePool]);

  return (
    <Select value={activePool?.id} onValueChange={handleChange} defaultValue={activePool?.id}>
      <SelectTrigger className="w-full">
        <SelectValue placeholder="Select a fruit" />
      </SelectTrigger>
      <SelectContent>
        <SelectGroup>
          <SelectLabel>Pool</SelectLabel>
          {pools.map((pool) => (
            <SelectItem key={pool.id} value={pool.id}>
              {pool.name}
            </SelectItem>
          ))}
        </SelectGroup>
      </SelectContent>
    </Select>
  )
}