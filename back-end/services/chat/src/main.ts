import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { ConfigService } from '@nestjs/config';
import { Logger } from '@nestjs/common';
import { NatsOptions, Transport } from '@nestjs/microservices';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  const config = app.get(ConfigService);

  app.connectMicroservice<NatsOptions>({
    transport: Transport.NATS,
    options: {
      servers: [config.get('natsUri')],
    },
  });

  await app.startAllMicroservices();
  Logger.log(`NATS Microservice is listening`);
  const port = config.get('port');
  await app.listen(port);
  Logger.log(`Listening on port ${port}`);
}
bootstrap();
