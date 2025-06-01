'use client';

import { useEffect, useState } from 'react';
import { Dream } from '@/lib/types/dream';
import { DreamService } from '@/lib/services/dream-service';
import Link from 'next/link';
import { ChevronRightIcon } from './icons/ChevronRightIcon';

export default function DreamList() {
  const dreamService = new DreamService();
  const [dreams, setDreams] = useState<Dream[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadDreams = async () => {
      try {
        const fetchedDreams = await dreamService.getAll();
        setDreams(fetchedDreams);
      } catch (err) {
        setError('Failed to load dreams');
        console.error('Error loading dreams:', err);
      }
    };

    loadDreams();
  }, []);

  if (error) {
    return (
      <div className="p-4 bg-red-900/50 text-red-200 rounded">
        {error}
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <h2 className="text-2xl font-semibold mb-4 text-purple-300">Your Dreams</h2>
      {dreams.length === 0 ? (
        <p className="text-gray-500">No dreams recorded yet.</p>
      ) : (
        <div className="space-y-4">
          {dreams.map((dream) => (
            <Link
              key={dream.id}
              href={`/dream/${dream.id}`}
              className="group relative block p-4 border border-gray-700 rounded hover:bg-gray-800/50 transition-all cursor-pointer"
            >
              <div className="pr-8">
                <p className="text-gray-200 whitespace-pre-wrap line-clamp-3">{dream.dream}</p>
                <p className="text-sm text-gray-400 mt-2">
                  {new Date(dream.created_at).toLocaleDateString()}
                </p>
              </div>
              <div className="absolute right-4 top-1/2 -translate-y-1/2">
                <ChevronRightIcon className="w-5 h-5 text-gray-500 opacity-0 group-hover:opacity-100 group-hover:text-purple-400 transition-all duration-200" />
              </div>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
