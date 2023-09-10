import { Module } from '@nestjs/common';
import { ChatGateway } from './chat.gateway';
import { DatabaseModule } from 'src/database/database.module';
import { ChatService } from './chat.service';
import { messageProviders } from './message.providers';
import { ChatController } from './chat.controller';

@Module({
  imports: [DatabaseModule],
  providers: [...messageProviders, ChatGateway, ChatService],
  controllers: [ChatController],
})
export class ChatModule {}
