'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { Dream } from '@/lib/types/dream';
import { DreamService } from '@/lib/services/dream-service';
import Link from 'next/link';
import { PencilIcon } from '@/app/components/icons/PencilIcon';
import { XMarkIcon } from '@/app/components/icons/XMarkIcon';
import { CheckIcon } from '@/app/components/icons/CheckIcon';
import { TrashIcon } from '@/app/components/icons/TrashIcon';

export default function DreamDetails() {
  const params = useParams();
  const router = useRouter();
  const dreamId = params.id as string;
  const dreamService = new DreamService();
  const [dream, setDream] = useState<Dream | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [isGeneratingImage, setIsGeneratingImage] = useState(false);
  const [isEditing, setIsEditing] = useState(false);
  const [editedDream, setEditedDream] = useState<string>('');

  useEffect(() => {
    loadDream();
  }, [dreamId]);

  const loadDream = async () => {
    try {
      const fetchedDream = await dreamService.getById(dreamId);
      setDream(fetchedDream);
      setEditedDream(fetchedDream.dream);
    } catch (err) {
      setError('Failed to load dream');
      console.error('Error loading dream:', err);
    }
  };

  const handleGenerateImage = async () => {
    setIsGeneratingImage(true);
    try {
      const updatedDream = await dreamService.generateImage(dreamId);
      setDream(updatedDream);
    } catch (err) {
      setError('Failed to generate image');
      console.error('Error generating image:', err);
    } finally {
      setIsGeneratingImage(false);
    }
  };

  const handleEdit = () => {
    setIsEditing(true);
  };

  const handleSaveEdit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!dream) return;
    
    try {
      const updatedDream = await dreamService.update(dream.id, editedDream);
      setDream(updatedDream);
      setIsEditing(false);
    } catch (err) {
      setError('Failed to update dream');
      console.error('Error updating dream:', err);
    }
  };

  const handleCancelEdit = () => {
    if (!dream) return;
    setEditedDream(dream.dream);
    setIsEditing(false);
  };

  const handleDelete = async () => {
    if (!dream || !confirm('Are you sure you want to delete this dream?')) return;

    try {
      const result = await dreamService.removeDelete(dream.id);
      router.replace('/');
    } catch (err) {
      setError('Failed to delete dream');
      console.error('Error deleting dream:', err);
    }
  };

  if (error) {
    return (
      <div className="min-h-screen p-8 bg-gray-900 text-gray-100">
        <div className="p-4 bg-red-900/50 text-red-200 rounded relative">
          <div className="pr-8">{error}</div>
          <button
            onClick={() => setError(null)}
            className="absolute right-2 top-1/2 -translate-y-1/2 p-1 text-red-200 hover:text-red-100 rounded-full hover:bg-red-800/30 transition-colors"
            title="Dismiss"
          >
            <XMarkIcon className="w-5 h-5" />
          </button>
        </div>
      </div>
    );
  }

  if (!dream) {
    return (
      <div className="min-h-screen p-8 bg-gray-900 text-gray-100">
        <p className="text-gray-400">Loading dream...</p>
      </div>
    );
  }

  return (
    <main className="min-h-screen p-8 bg-gray-900 text-gray-100">
      <div className="max-w-2xl mx-auto">
        <div className="flex items-center justify-between mb-8">
          <div className="flex items-center gap-4">
            <Link
              href="/"
              className="text-purple-400 hover:text-purple-300 transition-colors"
            >
              ← Back to Dreams
            </Link>
            <h1 className="text-4xl font-bold text-purple-300">Dream Details</h1>
          </div>
          <div className="flex gap-2">
            <button
              onClick={handleEdit}
              className="p-2 text-gray-400 hover:text-gray-300 rounded-full hover:bg-gray-700/30 transition-colors"
              title="Edit"
            >
              <PencilIcon className="w-5 h-5" />
            </button>
            <button
              onClick={handleDelete}
              className="p-2 text-gray-400 hover:text-red-400 rounded-full hover:bg-gray-700/30 transition-colors"
              title="Delete"
            >
              <TrashIcon className="w-5 h-5" />
            </button>
          </div>
        </div>
        
        <div className="bg-gray-800/50 border border-gray-700 rounded-lg p-6 mb-6">
          {isEditing ? (
            <form onSubmit={handleSaveEdit} className="relative">
              <textarea
                value={editedDream}
                onChange={(e) => setEditedDream(e.target.value)}
                className="w-full p-2 rounded resize-none focus:outline-none focus:ring-1 focus:ring-purple-500 bg-gray-800 text-gray-100 border-gray-700"
                rows={5}
                autoFocus
              />
              <div className="absolute right-2 top-2 flex gap-2">
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
            <>
              <p className="text-gray-200 text-lg mb-4 whitespace-pre-wrap">{dream.dream}</p>
              <p className="text-sm text-gray-400">
                Recorded on {new Date(dream.created_at).toLocaleDateString()}
              </p>
            </>
          )}
        </div>

        <div className="bg-gray-800/50 border border-gray-700 rounded-lg p-6">
          <h2 className="text-xl font-semibold mb-4 text-purple-300">Dream Visualization</h2>
          {dream.image_url ? (
            <div className="aspect-square w-full relative mb-4 rounded-lg overflow-hidden">
              <img
                src={dream.image_url}
                alt="AI generated visualization of the dream"
                className="object-cover w-full h-full"
              />
            </div>
          ) : (
            <div className="text-center p-8 bg-gray-800/30 rounded-lg mb-4">
              <p className="text-gray-400">No image generated yet</p>
            </div>
          )}
          
          <button
            onClick={handleGenerateImage}
            disabled={isGeneratingImage}
            className={`w-full py-2 px-4 rounded-lg font-medium transition-colors
              ${isGeneratingImage 
                ? 'bg-purple-700/50 cursor-not-allowed' 
                : 'bg-purple-600 hover:bg-purple-700'
              }`}
          >
            {isGeneratingImage ? 'Generating...' : 'Generate Image'}
          </button>
        </div>
      </div>
    </main>
  );
} 