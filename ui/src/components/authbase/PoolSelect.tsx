"use client";

import {createListCollection, ListCollection} from "@chakra-ui/react";
import {useEffect, useState} from "react";
import {usePoolStore} from "../../store/pool.ts";
import {
  SelectContent,
  SelectItem,
  SelectRoot,
  SelectTrigger,
  SelectValueText,
} from "../ui/select";

export const PoolSelect = () => {
  const pools = usePoolStore((state) => state.pools);
  const activePool = usePoolStore((state) => state.activePool);
  const setActivePool = usePoolStore((state) => state.setActivePool);
  const [collections, setCollections] = useState<ListCollection<any>>(createListCollection({items: []}));

  useEffect(() => {
    if (pools) {
      setCollections(
        createListCollection({
          items: pools.map((pool) => ({
            label: pool.name,
            value: pool.id,
          })),
        })
      );
      setActivePool(pools[0] || null);
    }
  }, [pools]);

  return (
    <SelectRoot
      collection={collections} size="sm" width="full" border={'1px solid #aaa'} borderRadius={6}
      value={activePool ? [activePool.id] : []}
      onValueChange={({value}) => {
        const pool = pools?.find((pool) => pool.id === value[0]);
        setActivePool(pool || null);
      }}>
      <SelectTrigger>
        <SelectValueText placeholder="Select Pool"/>
      </SelectTrigger>
      <SelectContent>
        {collections.items.map((movie) => (
          <SelectItem item={movie} key={movie.value} p={2}>
            {movie.label}
          </SelectItem>
        ))}
      </SelectContent>
    </SelectRoot>
  );
};
