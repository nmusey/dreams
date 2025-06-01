'use client';

import { useState, useEffect } from 'react';
import { useRouter, useParams } from 'next/navigation';
import Link from 'next/link';
import { Dream } from '@/lib/types/dream';
import { DreamService } from '@/lib/services/dream-service';

// Icons
import { PencilIcon, TrashIcon, XMarkIcon, CheckIcon } from '@heroicons/react/24/outline';

type GenerationStatus = {
  isGenerating: boolean;
  message: string;
  position?: number;
};

type ImageGenerationResponse = {
  message: string;
  position?: number;
  dreamId: string | number;
};

export default function DreamDetails() {
  const params = useParams();
  const router = useRouter();
  const dreamId = params.id as string;
  const dreamService = new DreamService();

  const [dream, setDream] = useState<Dream | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [generationStatus, setGenerationStatus] = useState<GenerationStatus>({
    isGenerating: false,
    message: '',
    position: undefined
  });
  const [isEditing, setIsEditing] = useState(false);
  const [editedDream, setEditedDream] = useState('');

  const updateGenerationStatus = (status: { position?: number; message: string; isGenerating?: boolean }) => {
    setGenerationStatus({
      isGenerating: status.isGenerating ?? true,
      message: status.message,
      position: status.position
    });
  };

  // Load the dream data when the component mounts
  const loadDream = async () => {
    if (!dreamId) return;
    
    setIsLoading(true);
    setError(null);
    
    try {
      const fetchedDream = await dreamService.getById(dreamId);
      if (fetchedDream) {
        setDream(fetchedDream);
        setEditedDream(fetchedDream.dream);
      } else {
        setError('Dream not found');
      }
    } catch (err) {
      console.error('Error loading dream:', err);
      setError('Failed to load dream. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  // Handle image generation
  const handleGenerateImage = async () => {
    if (!dream) return;
    
    try {
      setError(null);
      setGenerationStatus({
        isGenerating: true,
        message: 'Preparing your dream for generation...',
        position: undefined
      });
      
      // Convert dream.id to string for the API call
      const dreamId = typeof dream.id === 'number' ? dream.id.toString() : dream.id;
      
      // The generateImage method should handle the string ID
      const response = await dreamService.generateImage(dreamId);
      
      // Handle the response according to its type
      const status: ImageGenerationResponse = {
        message: response.message || 'Image generation in progress',
        position: 'position' in response ? response.position : undefined,
        dreamId: dreamId
      };
      
      setGenerationStatus({
        isGenerating: true,
        message: status.message,
        position: status.position
      });
    } catch (err) {
      console.error('Error generating image:', err);
      setError('Failed to generate image. Please try again.');
      setGenerationStatus({
        isGenerating: false,
        message: 'Generation failed',
        position: undefined
      });
    }
  };

  // Handle edit mode
  const handleEdit = () => {
    setIsEditing(true);
  };

  // Handle save edit
  const handleSaveEdit = async () => {
    if (!dream || !editedDream.trim()) return;
    
    try {
      // Pass the edited dream content as a string directly
      const updatedDream = await dreamService.update(dream.id, editedDream);
      setDream(updatedDream);
      setIsEditing(false);
    } catch (err) {
      console.error('Error updating dream:', err);
      setError('Failed to update dream. Please try again.');
    }
  };

  // Handle cancel edit
  const handleCancelEdit = () => {
    setIsEditing(false);
    if (dream) {
      setEditedDream(dream.dream);
    }
  };

  // Handle delete
  const handleDelete = async () => {
    if (!dream) return;
    
    try {
      // Use removeDelete with string ID
      await dreamService.removeDelete(dream.id.toString());
      router.push('/dreams');
    } catch (err) {
      console.error('Error deleting dream:', err);
      setError('Failed to delete dream. Please try again.');
    }
  };

  // Initial load
  useEffect(() => {
    loadDream();
  }, [dreamId]);

  // Loading state
  if (isLoading) {
    return (
      <div className="min-h-screen p-8 bg-gray-900 text-gray-100">
        <div className="max-w-3xl mx-auto">
          <div className="bg-gray-800 rounded-lg p-6 text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-500 mx-auto mb-4"></div>
            <p className="text-gray-400">Loading dream...</p>
          </div>
        </div>
      </div>
    );
  }

  // Error state
  if (error) {
    return (
      <div className="min-h-screen p-8 bg-gray-900 text-gray-100">
        <div className="max-w-3xl mx-auto">
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
          <Link href="/" className="text-purple-300 hover:text-purple-200">
            &larr; Back to dreams
          </Link>
        </div>
      </div>
    );
  }

  // No dream data
  if (!dream) {
    return (
      <div className="min-h-screen p-8 bg-gray-900 text-gray-100">
        <div className="max-w-3xl mx-auto">
          <div className="bg-gray-800 rounded-lg p-6 text-center">
            <p className="text-gray-400">No dream data available</p>
            <button
              onClick={loadDream}
              className="mt-4 px-4 py-2 bg-purple-600 text-white rounded hover:bg-purple-700 transition-colors"
            >
              Reload
            </button>
          </div>
        </div>
      </div>
    );
  }

  // Main content
  return (
    <div className="min-h-screen p-8 bg-gray-900 text-gray-100">
      <div className="max-w-3xl mx-auto">
        <div className="flex justify-between items-center mb-6">
          <Link href="/" className="text-purple-300 hover:text-purple-200">
            &larr; Back to dreams
          </Link>
          <div className="space-x-2">
            {isEditing ? (
              <>
                <button
                  onClick={handleSaveEdit}
                  className="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 flex items-center transition-colors"
                >
                  <CheckIcon className="h-5 w-5 mr-1" />
                  Save
                </button>
                <button
                  onClick={handleCancelEdit}
                  className="px-4 py-2 bg-gray-700 text-gray-200 rounded hover:bg-gray-600 flex items-center transition-colors"
                >
                  <XMarkIcon className="h-5 w-5 mr-1" />
                  Cancel
                </button>
              </>
            ) : (
              <>
                <button
                  onClick={handleEdit}
                  className="px-4 py-2 bg-purple-600 text-white rounded hover:bg-purple-700 flex items-center transition-colors"
                >
                  <PencilIcon className="h-5 w-5 mr-1" />
                  Edit
                </button>
                <button
                  onClick={handleDelete}
                  className="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 flex items-center transition-colors"
                >
                  <TrashIcon className="h-5 w-5 mr-1" />
                  Delete
                </button>
              </>
            )}
          </div>
        </div>

        <div className="bg-gray-800 rounded-lg p-6">
          <div className="mb-6">
            <div className="text-sm text-gray-400 mb-4">
              Created on {new Date(dream.created_at).toLocaleDateString()}
            </div>
            
            {isEditing ? (
              <form onSubmit={handleSaveEdit}>
                <textarea
                  value={editedDream}
                  onChange={(e) => setEditedDream(e.target.value)}
                  className="w-full p-3 rounded bg-gray-700 text-gray-100 border-gray-600 focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                  rows={4}
                  autoFocus
                />
              </form>
            ) : (
              <p className="text-gray-200 whitespace-pre-line">
                {dream.dream}
              </p>
            )}
          </div>

          {dream.image_url ? (
            <div className="mt-8 border-t border-gray-700 pt-6">
              <h2 className="text-xl font-semibold mb-4 text-purple-300">Generated Image</h2>
              <div className="relative">
                <img
                  src={dream.image_url}
                  alt="Generated from dream"
                  className="w-full h-auto rounded-lg"
                  onError={(e) => {
                    console.error('Error loading image:', dream.image_url);
                    const target = e.target as HTMLImageElement;
                    target.onerror = null;
                    target.src = '/placeholder-image.png';
                  }}
                />
              </div>
            </div>
          ) : (
            <div className="mt-8 border-t border-gray-700 pt-6">
              <h2 className="text-xl font-semibold mb-4 text-purple-300">Generate Image</h2>
              <div className="bg-gray-700/50 p-6 rounded-lg border border-dashed border-gray-600 text-center">
                <p className="text-gray-300 mb-4">
                  {generationStatus.isGenerating
                    ? generationStatus.message
                    : 'Generate an image based on your dream'}
                </p>
                <button
                  onClick={handleGenerateImage}
                  disabled={generationStatus.isGenerating}
                  className={`px-6 py-3 rounded-lg font-medium flex items-center mx-auto ${
                    generationStatus.isGenerating
                      ? 'bg-gray-600 cursor-not-allowed'
                      : 'bg-gradient-to-r from-purple-600 to-blue-600 hover:from-purple-700 hover:to-blue-700 text-white shadow-md'
                  }`}
                >
                  {generationStatus.isGenerating ? (
                    <>
                      <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                      </svg>
                      Processing...
                    </>
                  ) : (
                    'Generate Image'
                  )}
                </button>
                {generationStatus.position !== undefined && generationStatus.position >= 0 && (
                  <p className="mt-3 text-sm text-gray-400">
                    Position in queue: {generationStatus.position + 1}
                  </p>
                )}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
