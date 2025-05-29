import { Entity, PrimaryGeneratedColumn, Column, EntityRepository, Repository } from 'typeorm';

@Entity()
export class Dream {
  @PrimaryGeneratedColumn()
  id: number;

  @Column({ type: 'text', nullable: false })
  dream: string;
}
