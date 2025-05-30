import React, { useEffect, useState } from 'react';
import fetchDreams from "@/app/components/DreamServerList";
import DreamClientList from "@/components/DreamClientList";

export default function DreamList() {
  const [dreams, setDreams] = useState<Dream[]>([]);

  useEffect(() => {
    const getDreams = async () => {
      try {
        const { dreams } = await fetchDreams();
        setDreams(dreams);
      } catch (error) {
        console.error('Error fetching dreams:', error);
      }
    };

    getDreams();
  }, []);

  return (
    <DreamClientList dreams={dreams} />
  );
}
