import { Entity, Column, PrimaryColumn } from 'typeorm';

@Entity('messages')
export class Message {
  @PrimaryColumn({ type: 'uuid', name: 'id' })
  id: string;

  @Column({ type: 'text', name: 'text' })
  text: string;

  @Column({ type: 'uuid', name: 'user_id' })
  userId: string;

  @Column({ type: 'timestamp with time zone', name: 'created_at' })
  createdAt: Date;

  @Column({ type: 'timestamp with time zone', name: 'updated_at' })
  updatedAt: Date;
}
