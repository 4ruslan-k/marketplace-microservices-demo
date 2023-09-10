import { Logger, UnauthorizedException } from '@nestjs/common';
import {
  OnGatewayConnection,
  OnGatewayDisconnect,
  SubscribeMessage,
  WebSocketGateway,
  WebSocketServer,
} from '@nestjs/websockets';
import { v4 as uuidv4 } from 'uuid';

import { Socket, Server } from 'socket.io';
import { ChatService } from './chat.service';

class AuthInfo {
  user_id: string;
  ip: string;
  user_agent: string;
  session_id: string;
}

@WebSocketGateway(80, {
  path: '/chat/socket.io',
  credentials: true,
  cors: {
    origin: ['http://localhost:3500'],
    credentials: true,
  },
})
export class ChatGateway implements OnGatewayConnection, OnGatewayDisconnect {
  @WebSocketServer()
  server: Server;

  private logger = new Logger('WebSocketServer');

  constructor(private readonly chatService: ChatService) {}

  async handleConnection(socket: Socket) {
    try {
      this.logger.log(`handleConnection -> Client connected: ${socket.id}`);
      const authHeader = socket.handshake.headers['x-authentication-info'];
      const authInfo = JSON.parse(String(authHeader));

      if (!authHeader || !authInfo) {
        return this.disconnect(socket);
      } else {
        socket.data.authInfo = authInfo;
        return this.server.to(socket.id);
      }
    } catch {
      return this.disconnect(socket);
    }
  }

  async handleDisconnect(socket: Socket) {
    // remove connection from DB
    this.logger.log(`handleDisconnect -> Client disconnected: ${socket.id}`);
    socket.disconnect();
  }

  private disconnect(socket: Socket) {
    this.logger.log(`disconnect -> Client disconnected: ${socket.id}`);

    socket.emit('Error', new UnauthorizedException());
    socket.disconnect();
  }

  @SubscribeMessage('addMessage')
  async onAddMessage(socket: Socket, message: any) {
    this.logger.log(`addMessage -> Client: ${socket.id} Message: ${message}`);
    const authInfo: AuthInfo = socket.data.authInfo;
    const addedMessage = await this.chatService.sendMessage({
      id: uuidv4(),
      text: message,
      userId: authInfo.user_id,
    });
    this.server.emit('newMessage', addedMessage);
  }
}
