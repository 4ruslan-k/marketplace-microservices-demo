import { NestFactory } from '@nestjs/core';
import { AppModule } from './app.module';
import { ConfigService } from '@nestjs/config';
import { Logger } from '@nestjs/common';
import { MicroserviceOptions, Transport } from '@nestjs/microservices';

async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  const config = app.get(ConfigService);

  const natsApp = NestFactory.createMicroservice<MicroserviceOptions>(
    AppModule,
    {
      transport: Transport.NATS,
      options: {
        servers: ['nats://localhost:4223'],
        url: 'nats://localhost:4223',
      },
    },
  );
  app.connectMicroservice(natsApp);
  // await natsApp.listen;
  const port = config.get('port');
  await app.listen(port);
  Logger.log(`Listening on port ${port}`);

  Logger.log(`NATS Microservice is listening`);
}
bootstrap();
