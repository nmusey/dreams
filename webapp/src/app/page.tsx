'use client';

import { useEffect, useState } from 'react';
import { Dream } from '@/lib/types/dream';
import { DreamService } from '@/lib/services/dream-service';
import { PencilIcon } from './components/icons/PencilIcon';
import { XMarkIcon } from './components/icons/XMarkIcon';
import { CheckIcon } from './components/icons/CheckIcon';
import { TrashIcon } from './components/icons/TrashIcon';

const dreamService = DreamService.getInstance();

export default function Home() {
  const [dreams, setDreams] = useState<Dream[]>([]);
  const [newDream, setNewDream] = useState('');
  const [editingDream, setEditingDream] = useState<Dream | null>(null);
  const [originalDream, setOriginalDream] = useState<Dream | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadDreams();
  }, []);

  const loadDreams = async () => {
    try {
      console.log('Loading dreams...');
      const fetchedDreams = await dreamService.getAll();
      console.log('Fetched dreams:', fetchedDreams);
      setDreams(fetchedDreams);
      setEditingDream(null);
      setOriginalDream(null);
      console.log('Dreams loaded successfully');
    } catch (err) {
      console.error('Error loading dreams:', err);
      setError('Failed to load dreams');
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

  const handleEdit = (dream: Dream) => {
    setOriginalDream(dream);
    setEditingDream({ ...dream });
  };

  const handleSaveEdit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingDream) return;
    if (editingDream.id === undefined || editingDream.id === null) {
      setError('Dream ID is missing. Cannot update.');
      return;
    }
    try {
      console.log('Saving dream:', editingDream);
      const updatedDream = await dreamService.update(editingDream.id, editingDream.dream);
      console.log('Dream updated successfully:', updatedDream);
      await loadDreams();
      console.log('Dreams reloaded after update');
      setEditingDream(null);
      setOriginalDream(null);
    } catch (err) {
      console.error('Error updating dream:', err);
      setError('Failed to update dream');
    }
  };

  const handleCancelEdit = () => {
    setEditingDream(null);
    setOriginalDream(null);
  };

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this dream?')) return;

    try {
      await dreamService.removeDelete(id);
      setDreams(dreams.filter(dream => dream.id !== id));
      setEditingDream(null);
      setOriginalDream(null);
    } catch (err) {
      setError('Failed to delete dream');
      console.error('Error deleting dream:', err);
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

      <div className="space-y-4">
        {dreams.map((dream) => (
          <div
            key={dream.id}
            className="p-4 border border-gray-700 rounded hover:bg-gray-800/50 relative group transition-colors"
          >
            <div className="flex items-start gap-4">
              <div className="flex-1 min-w-0">
                {editingDream && (editingDream?.id === dream.id) ? (
                  <form onSubmit={handleSaveEdit} className="relative">
                    <textarea
                      value={editingDream.dream}
                      onChange={(e) => setEditingDream({ ...editingDream, dream: e.target.value })}
                      className="w-full p-2 rounded resize-none focus:outline-none focus:ring-1 focus:ring-purple-500 bg-gray-800 text-gray-100 border-gray-700"
                      rows={3}
                      autoFocus
                    />
                    <div className="absolute right-0 top-0 flex gap-2">
                      <button
                        type="submit"
                        className="p-1 text-purple-400 hover:text-purple-300 rounded-full hover:bg-purple-900/30 transition-colors"
                        title="Save"
                      >
                        <CheckIcon className="w-5 h-5" />
                      </button>
                      <button
                        type="button"
                        onClick={handleCancelEdit}
                        className="p-1 text-gray-400 hover:text-gray-300 rounded-full hover:bg-gray-700/30 transition-colors"
                        title="Cancel"
                      >
                        <XMarkIcon className="w-5 h-5" />
                      </button>
                    </div>
                  </form>
                ) : (
                  <div className="relative">
                    <p className="text-gray-200 whitespace-pre-wrap">{dream.dream}</p>
                    <p className="text-sm text-gray-400 mt-2">
                      {formatDate(dream.created_at)}
                    </p>
                    <div className="absolute right-0 top-0 flex gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                      <button
                        onClick={() => handleEdit(dream)}
                        className="p-1 text-gray-400 hover:text-gray-300 rounded-full hover:bg-gray-700/30 transition-colors"
                        title="Edit"
                      >
                        <PencilIcon className="w-5 h-5" />
                      </button>
                      <button
                        onClick={() => handleDelete(dream.id)}
                        className="p-1 text-gray-400 hover:text-red-400 rounded-full hover:bg-gray-700/30 transition-colors"
                        title="Delete"
                      >
                        <TrashIcon className="w-5 h-5" />
                      </button>
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>
        ))}
      </div>
    </main>
  );
}
