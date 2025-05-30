"use server";

import DreamRepository from "@/lib/repositories/dream-repository";
import { Dream } from "@/lib/entities/Dream"; 

interface FetchDreamsArgs {
  limit?: number;
  cursor?: string | null;
}

export async function fetchDreams(args: FetchDreamsArgs = {}): Promise<{ dreams: Dream[]; nextCursor: string | null }> {
  const dreamRepository = new DreamRepository();
  try {
    return await dreamRepository.findAll(args);
  } catch (error) {
    console.error('Error fetching dreams:', error);
    throw error;
  }
}
