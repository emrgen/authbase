import {create} from 'zustand'

interface Project {
  id: string
  name: string
}

interface ProjectStore {
  listProjectState: 'loading' | 'error' | 'success' | 'idle';
  projects: Project[];
  activeProject: Project | null;
  setProjects: (projects: Project[]) => void
  setActiveProject: (project: Project | null) => void
  setListProjectState: (state: 'loading' | 'error' | 'success') => void
}

export const useProjectStore = create<ProjectStore>((set) => ({
  listProjectState: 'idle',
  projects: [],
  activeProject: null,
  setProjects: (projects: Project[]) => {
    // Set the projects in the store
    set({projects})
    if (projects.length > 0) {
      // Set the first project as the active project
      set({activeProject: projects[0]})
    }
  },
  setActiveProject: (project: Project | null) => {
    // Set the active project in the store
    set({activeProject: project})
  },
  setListProjectState: (state: 'loading' | 'error' | 'success') => {
    // Set the list project state in the store
    set({listProjectState: state})
  },
}))