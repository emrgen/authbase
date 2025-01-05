"use client";

import { createListCollection } from "@chakra-ui/react";
import {
  SelectContent,
  SelectItem,
  SelectRoot,
  SelectTrigger,
  SelectValueText,
} from "../ui/select";

export const ProjectSelect = () => {
  return (
    <SelectRoot collection={frameworks} size="sm" width="full" border={'1px solid #aaa'} borderRadius={6}>
      <SelectTrigger>
        <SelectValueText placeholder="Select Project" />
      </SelectTrigger>
      <SelectContent>
        {frameworks.items.map((movie) => (
          <SelectItem item={movie} key={movie.value} p={2}>
            {movie.label}
          </SelectItem>
        ))}
      </SelectContent>
    </SelectRoot>
  );
};

const frameworks = createListCollection({
  items: [
    { label: 'master', value: 'master' },
    { label: 'dev', value: 'dev' },
    { label: 'feature', value: 'feature' },
    { label: 'bugfix', value: 'bugfix' },
  ],
});
