'use client';

import { useEffect, useState } from 'react';
import { Dream } from '@/lib/types/dream';
import { DreamService } from '@/lib/services/dream-service';

export default function DreamList() {
  const dreamService = new DreamService();
  const [dreams, setDreams] = useState<Dream[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadDreams();
  }, []);

  const loadDreams = async () => {
    try {
      const fetchedDreams = await dreamService.getAll();
      setDreams(fetchedDreams);
    } catch (err) {
      setError('Failed to load dreams');
      console.error('Error loading dreams:', err);
    }
  };

  if (error) {
    return (
      <div className="p-4 bg-red-100 text-red-700 rounded">
        {error}
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <h2 className="text-2xl font-semibold mb-4">Your Dreams</h2>
      {dreams.length === 0 ? (
        <p className="text-gray-500">No dreams recorded yet.</p>
      ) : (
        <div className="space-y-4">
          {dreams.map((dream) => (
            <div
              key={dream.id}
              className="p-4 border rounded hover:bg-gray-50 transition-colors"
            >
              <p className="text-gray-800">{dream.dream}</p>
              <p className="text-sm text-gray-500 mt-2">
                {new Date(dream.created_at).toLocaleDateString()}
              </p>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
