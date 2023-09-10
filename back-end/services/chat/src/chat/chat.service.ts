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
    const messageInput = {
      id: message.id,
      text: message.text,
      userId: message.userId,
      createdAt: new Date(),
      updatedAt: null,
    };
    await this.messageRepository.insert(messageInput);
    return messageInput;
  }

  async getMessages(): Promise<Message[]> {
    return this.messageRepository.find();
  }
}

export class CreateMessageDto {
  id: string;
  text: string;
  userId: string;
}
