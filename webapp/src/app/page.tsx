'use client';

import { useCallback, useRef, useState } from 'react';
import { DreamService } from '@/lib/services/dream-service';
import { XMarkIcon } from './components/icons/XMarkIcon';
import DreamList from './components/DreamList';

const dreamService = DreamService.getInstance();

export default function Home() {
  const [newDream, setNewDream] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [refreshKey, setRefreshKey] = useState(0);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!newDream.trim()) return;

    try {
      await dreamService.create(newDream);
      setNewDream('');
      setRefreshKey(key => key + 1); // Trigger a refresh
    } catch (err) {
      setError('Failed to create dream');
      console.error('Error creating dream:', err);
    }
  };

  return (
    <main className="min-h-screen p-8 bg-gray-900 text-gray-100">
      <h1 className="text-4xl font-bold mb-8 text-purple-300">Dream Journal</h1>
      
      <form onSubmit={handleSubmit} className="mb-8">
        <div className="flex gap-4">
          <input
            type="text"
            value={newDream}
            onChange={(e) => setNewDream(e.target.value)}
            placeholder="Enter your dream..."
            className="flex-1 p-2 rounded bg-gray-800 border-gray-700 text-gray-100 placeholder-gray-500 focus:outline-none focus:ring-1 focus:ring-purple-500 focus:border-purple-500"
          />
          <button
            type="submit"
            className="px-4 py-2 bg-purple-600 text-white rounded hover:bg-purple-700 transition-colors"
          >
            Save Dream
          </button>
        </div>
      </form>

      {error && (
        <div className="mb-4 p-4 bg-red-900/50 text-red-200 rounded relative group">
          <div className="pr-8">{error}</div>
          <button
            onClick={() => setError(null)}
            className="absolute right-2 top-1/2 -translate-y-1/2 p-1 text-red-200 hover:text-red-100 rounded-full hover:bg-red-800/30 transition-colors"
            title="Dismiss"
          >
            <XMarkIcon className="w-5 h-5" />
          </button>
        </div>
      )}

      <DreamList key={refreshKey} />
    </main>
  );
}
