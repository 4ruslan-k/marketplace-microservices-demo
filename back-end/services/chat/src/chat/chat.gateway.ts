import { Logger, UnauthorizedException } from '@nestjs/common';
import {
  OnGatewayConnection,
  OnGatewayDisconnect,
  SubscribeMessage,
  WebSocketGateway,
  WebSocketServer,
} from '@nestjs/websockets';

import { Socket, Server } from 'socket.io';

@WebSocketGateway(80, {
  path: '/chat/socket.io',
  cors: {
    origin: ['http://localhost:3500'],
  },
})
export class ChatGateway implements OnGatewayConnection, OnGatewayDisconnect {
  @WebSocketServer()
  server: Server;
  private logger = new Logger('WebSocketServer');

  constructor() {}

  async handleConnection(socket: Socket) {
    try {
      this.logger.log(`handleConnection -> Client connected: ${socket.id}`);

      const user = {
        id: 1,
        name: 'test',
      };
      if (!user) {
        return this.disconnect(socket);
      } else {
        socket.data.user = user;
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
    socket.emit('newMessage', `Client ${socket.id} sent message: ${message}`);
  }
}