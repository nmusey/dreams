'use client';

import { useEffect, useState } from 'react';
import { Dream } from '@/lib/types/dream';
import { DreamService } from '@/lib/services/dream-service';

const dreamService = DreamService.getInstance();

export default function Home() {
  const [dreams, setDreams] = useState<Dream[]>([]);
  const [newDream, setNewDream] = useState('');
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadDreams();
  }, []);

  const loadDreams = async () => {
    try {
      const fetchedDreams = await dreamService.getAll();
      console.log('Fetched dreams:', fetchedDreams);
      setDreams(fetchedDreams);
    } catch (err) {
      setError('Failed to load dreams');
      console.error('Error loading dreams:', err);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newDream.trim()) return;

    try {
      await dreamService.create(newDream);
      setNewDream('');
      await loadDreams(); // Reload dreams after creating a new one
    } catch (err) {
      setError('Failed to create dream');
      console.error('Error creating dream:', err);
    }
  };

  const formatDate = (dateString: string) => {
    try {
      const date = new Date(dateString);
      if (isNaN(date.getTime())) {
        console.error('Invalid date:', dateString);
        return 'Invalid date';
      }
      return date.toLocaleDateString('en-US', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      });
    } catch (err) {
      console.error('Error formatting date:', err);
      return 'Invalid date';
    }
  };

  return (
    <main className="min-h-screen p-8">
      <h1 className="text-4xl font-bold mb-8">Dream Journal</h1>
      
      <form onSubmit={handleSubmit} className="mb-8">
        <div className="flex gap-4">
          <input
            type="text"
            value={newDream}
            onChange={(e) => setNewDream(e.target.value)}
            placeholder="Enter your dream..."
            className="flex-1 p-2 border rounded"
          />
          <button
            type="submit"
            className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
          >
            Save Dream
          </button>
        </div>
      </form>

      {error && (
        <div className="mb-4 p-4 bg-red-100 text-red-700 rounded">
          {error}
        </div>
      )}

      <div className="space-y-4">
        {dreams.map((dream) => (
          <div
            key={dream.id}
            className="p-4 border rounded hover:bg-gray-50"
          >
            <p className="text-gray-800">{dream.dream}</p>
            <p className="text-sm text-gray-500 mt-2">
              {formatDate(dream.created_at)}
            </p>
          </div>
        ))}
      </div>
    </main>
  );
}
