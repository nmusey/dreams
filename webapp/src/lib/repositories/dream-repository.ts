import { getRepository } from 'typeorm';
import { Dream } from '../entities/Dream';

export class DreamRepository {
  private dreamRepository = getRepository(Dream);

    async create(dream: string): Promise<void> {
        // TODO - fix this
        // return this.dreamRepository.create();
    }

    async findAll(): Promise<Dream[]> {
        return this.dreamRepository.find();
    }
}
