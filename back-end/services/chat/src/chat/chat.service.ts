import { Injectable, Inject } from '@nestjs/common';
import { Repository } from 'typeorm';
import { MESSAGE_REPOSITORY } from './message.providers';
import { Message } from './message.entity';

@Injectable()
export class ChatService {
  constructor(
    @Inject(MESSAGE_REPOSITORY)
    private messageRepository: Repository<Message>,
  ) {}

  async sendMessage(message: CreateMessageDto) {
    await this.messageRepository.insert({
      id: message.id,
      text: message.text,
      userId: message.userId,
      createdAt: new Date(),
      updatedAt: null,
    });
  }
}

export class CreateMessageDto {
  id: string;
  text: string;
  userId: string;
}
