import { Module } from '@nestjs/common';
import { ChatGateway } from './chat.gateway';
import { DatabaseModule } from 'src/database/database.module';
import { ChatService } from './chat.service';
import { messageProviders } from './message.providers';

@Module({
  imports: [DatabaseModule],
  providers: [...messageProviders, ChatGateway, ChatService],
})
export class ChatModule {}
