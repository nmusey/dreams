"use client"

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import DreamRepository from '@/lib/repositories/dream-repository';

export default function DreamList() {
  const [dreams, setDreams] = useState<Dream[]>([]);
  const router = useRouter();
  const dreamRepository = new DreamRepository();

  useEffect(() => {
    const fetchDreams = async () => {
      try {
        const fetchedDreams = await dreamRepository.findAll();
        setDreams(fetchedDreams);
      } catch (error) {
        console.error('Error fetching dreams:', error);
      }
    };

    fetchDreams();
  }, []);

  return (
    <div>
      <h2>Current Dreams</h2>
      <ul>
        {dreams.map((dream) => (
          <li key={dream.id}>{dream.dream}</li>
        ))}
      </ul>
    </div>
  );
}
