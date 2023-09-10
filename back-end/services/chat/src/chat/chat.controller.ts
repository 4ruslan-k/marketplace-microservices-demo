import { Controller, Get, Logger } from '@nestjs/common';
import { ChatService } from './chat.service';

@Controller('/v1/chat/messages')
export class ChatController {
  private logger = new Logger('UserControllers');
  constructor(private readonly chatService: ChatService) {}

  @Get()
  async getMessages() {
    const messages = await this.chatService.getMessages();
    const messagesOutput = messages.map((message) => ({
      id: message.id,
      text: message.text,
      type: 'message',
      createdAt: message.createdAt,
    }));
    return { type: 'list', items: messagesOutput };
  }
}
