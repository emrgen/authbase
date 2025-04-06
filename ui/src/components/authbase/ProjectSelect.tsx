"use client";

import {createListCollection, ListCollection} from "@chakra-ui/react";
import {useEffect, useState} from "react";
import {useProjectStore} from "../../store/project.ts";
import {
  SelectContent,
  SelectItem,
  SelectRoot,
  SelectTrigger,
  SelectValueText,
} from "../ui/select";


export const ProjectSelect = () => {
  const projects = useProjectStore((state) => state.projects);
  const activeProject = useProjectStore((state) => state.activeProject);
  // const listProjectState = useProjectStore((state) => state.listProjectState);
  const setActiveProject = useProjectStore((state) => state.setActiveProject);

  const [collections, setCollections] = useState<ListCollection<any>>(createListCollection({items: []}));

  useEffect(() => {
    if (projects) {
      setCollections(
        createListCollection({
          items: projects.map((project) => ({
            label: project.name,
            value: project.id,
          })),
        })
      );
    }
  }, [projects]);


  return (
    <SelectRoot
      collection={collections} size="sm" width="full" border={'1px solid #aaa'}
      borderRadius={6}
      value={activeProject ? [activeProject.id] : []}
      onValueChange={({value}) => {
        const project = projects?.find((project) => project.id === value[0]);
        if (!project) {
          console.error("Project not found");
          return;
        }
        setActiveProject(project);
      }}>
      <SelectTrigger>
        <SelectValueText placeholder="Select Project"/>
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
