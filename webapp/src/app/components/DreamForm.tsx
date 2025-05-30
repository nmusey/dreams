"use client"

import { useState } from 'react';
import { useRouter } from 'next/navigation';

export default function DreamForm() {
  const [dream, setDream] = useState('');
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      await fetch('/api/dreams', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ dream }),
      });
      router.push('/');
    } catch (error) {
      console.error('Error saving dream:', error);
    }
  };

  return (
    <form onSubmit={handleSubmit}>
      <div>
        <label htmlFor="dream">Enter your dream:</label>
        <textarea
          id="dream"
          value={dream}
          onChange={(e) => setDream(e.target.value)}
          required
        />
      </div>
      <button type="submit">Save Dream</button>
    </form>
  );
}
